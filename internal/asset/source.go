package asset

type Source struct {
	Name                        string
	Path                        string
	AdminEsbuildCompatible      bool
	StorefrontEsbuildCompatible bool
	DisableSass                 bool
	NpmStrict                   bool
}
