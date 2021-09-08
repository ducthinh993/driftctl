package terraform

import (
	"fmt"

	version2 "github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/tfdiags"
)

type LockFileProvider struct {
	Addr        addrs.Provider
	Constraints string   `hcl:"constraints"`
	Hashes      []string `hcl:"hashes"`
	Version     string   `hcl:"version"`
}

// Locks is the top-level type representing the information retained in a
// dependency lock file.
//
// Locks and the other types used within it are mutable via various setter
// methods, but they are not safe for concurrent  modifications, so it's the
// caller's responsibility to prevent concurrent writes and writes concurrent
// with reads.
type Locks struct {
	providers map[addrs.Provider]*LockFileProvider
}

// NewLocks constructs and returns a new Locks object that initially contains
// no locks at all.
func NewLocks() *Locks {
	return &Locks{
		providers: make(map[addrs.Provider]*LockFileProvider),
	}
}

// Provider returns the stored lock for the given provider, or nil if that
// provider currently has no lock.
func (l *Locks) Provider(addr *addrs.Provider) *LockFileProvider {
	return l.providers[*addr]
}

// ProviderIsLockable returns true if the given provider is eligible for
// version locking.
//
// Currently, all providers except builtin and legacy providers are eligible
// for locking.
func ProviderIsLockable(addr addrs.Provider) bool {
	return !(addr.IsBuiltIn() || addr.IsLegacy())
}

func LoadLocksFromFile(filename string) (*Locks, error) {
	return loadLocks(func(parser *hclparse.Parser) (*hcl.File, hcl.Diagnostics) {
		return parser.ParseHCLFile(filename)
	})
}

func loadLocks(loadParse func(*hclparse.Parser) (*hcl.File, hcl.Diagnostics)) (*Locks, error) {
	ret := NewLocks()

	var diags tfdiags.Diagnostics

	parser := hclparse.NewParser()
	f, hclDiags := loadParse(parser)
	diags = diags.Append(hclDiags)
	if f == nil {
		// If we encountered an error loading the file then those errors
		// should already be in diags from the above, but the file might
		// also be nil itself and so we can't decode from it.
		return ret, diags.Err()
	}

	moreDiags := decodeLocksFromHCL(ret, f.Body)
	diags = diags.Append(moreDiags)
	return ret, diags.Err()
}

func decodeLocksFromHCL(locks *Locks, body hcl.Body) error {
	var diags tfdiags.Diagnostics

	content, hclDiags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "provider",
				LabelNames: []string{"source_addr"},
			},
		},
	})
	diags = diags.Append(hclDiags)

	seenProviders := make(map[addrs.Provider]hcl.Range)
	for _, block := range content.Blocks {

		switch block.Type {
		case "provider":
			lock, moreDiags := decodeProviderLockFromHCL(block)
			diags = diags.Append(moreDiags)
			if lock == nil {
				continue
			}
			if previousRng, exists := seenProviders[lock.Addr]; exists {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate provider lock",
					Detail:   fmt.Sprintf("This lockfile already declared a lock for provider %s at %s.", lock.Addr.String(), previousRng.String()),
					Subject:  block.TypeRange.Ptr(),
				})
				continue
			}
			locks.providers[lock.Addr] = lock
			seenProviders[lock.Addr] = block.DefRange
		default:
			// Shouldn't get here because this should be exhaustive for
			// all of the block types in the schema above.
		}

	}

	return diags.Err()
}

func decodeProviderLockFromHCL(block *hcl.Block) (*LockFileProvider, error) {
	ret := &LockFileProvider{}
	var diags tfdiags.Diagnostics

	rawAddr := block.Labels[0]
	addr, moreDiags := addrs.ParseProviderSourceString(rawAddr)
	if moreDiags.HasErrors() {
		// The diagnostics from ParseProviderSourceString are, as the name
		// suggests, written with an intended audience of someone who is
		// writing a "source" attribute in a provider requirement, not
		// our lock file. Therefore we're using a less helpful, fixed error
		// here, which is non-ideal but hopefully okay for now because we
		// don't intend end-users to typically be hand-editing these anyway.
		return nil, fmt.Errorf("the provider source address for a provider lock must be a valid, fully-qualified address of the form \"hostname/namespace/type\"")
	}
	if !ProviderIsLockable(addr) {
		if addr.IsBuiltIn() {
			// A specialized error for built-in providers, because we have an
			// explicit explanation for why those are not allowed.
			return nil, fmt.Errorf("cannot lock a version for built-in provider %s. Built-in providers are bundled inside Terraform itself, so you can't select a version for them independently of the Terraform release you are currently running", addr)
		}
		// Otherwise, we'll use a generic error message.
		return nil, fmt.Errorf("provider source address %s is a special provider that is not eligible for dependency locking", addr)
	}
	if canonAddr := addr.String(); canonAddr != rawAddr {
		// We also require the provider addresses in the lock file to be
		// written in fully-qualified canonical form, so that it's totally
		// clear to a reader which provider each block relates to. Again,
		// we expect hand-editing of these to be atypical so it's reasonable
		// to be stricter in parsing these than we would be in the main
		// configuration.
		return nil, fmt.Errorf("the provider source address for this provider lock must be written as %q, the fully-qualified and normalized form", canonAddr)
	}

	ret.Addr = addr

	content, hclDiags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "version", Required: true},
			{Name: "constraints"},
			{Name: "hashes"},
		},
	})
	diags = diags.Append(hclDiags)

	ver, moreDiags := decodeProviderVersionArgument(addr, content.Attributes["version"])
	if ver != nil {
		ret.Version = ver.String()
	}
	diags = diags.Append(moreDiags)

	return ret, diags.Err()
}

func decodeProviderVersionArgument(provider addrs.Provider, attr *hcl.Attribute) (*version2.Version, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics
	if attr == nil {
		// It's not okay to omit this argument, but the caller should already
		// have generated diagnostics about that.
		return nil, diags
	}
	expr := attr.Expr

	var raw *string
	hclDiags := gohcl.DecodeExpression(expr, nil, &raw)
	diags = diags.Append(hclDiags)
	if hclDiags.HasErrors() {
		return nil, diags
	}
	if raw == nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Missing required argument",
			Detail:   "A provider lock block must contain a \"version\" argument.",
			Subject:  expr.Range().Ptr(), // the range for a missing argument's expression is the body's missing item range
		})
		return nil, diags
	}
	ver, err := version2.NewVersion(*raw)
	if err != nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider version number",
			Detail:   fmt.Sprintf("The selected version number for provider %s is invalid: %s.", provider, err),
			Subject:  expr.Range().Ptr(),
		})
	}
	if canon := ver.String(); canon != *raw {
		// Canonical forms are required in the lock file, to reduce the risk
		// that a file diff will show changes that are entirely cosmetic.
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider version number",
			Detail:   fmt.Sprintf("The selected version number for provider %s must be written in normalized form: %q.", provider, canon),
			Subject:  expr.Range().Ptr(),
		})
	}
	return ver, diags
}
