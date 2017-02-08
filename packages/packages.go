package packages

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
)

// Example Package xml:
//
// <Package>
//   <Name>golang</Name>
//   <Summary xml:lang="en">The Go programming language.</Summary>
//   <Description xml:lang="en">Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.</Description>
//   <PartOf>programming</PartOf>
//   <License>BSD-3-Clause</License>
//   <RuntimeDependencies>
//     <Dependency>glibc</Dependency>
//   </RuntimeDependencies>
//   <Replaces>
//     <Package>golang-binary</Package>
//   </Replaces>
//   <History>
//     <Update release="13">
//       <Date>2017-01-26</Date>
//       <Version>1.7.5</Version>
//       <Comment>Update to 1.7.5, drop tzdata patch since it is no longer needed for 1.7 branch.</Comment>
//       <Name>Joshua Strobl</Name>
//       <Email>joshua@stroblindustries.com</Email>
//     </Update>
//   </History>
//   <BuildHost>solus-build-server</BuildHost>
//   <Distribution>Solus</Distribution>
//   <DistributionRelease>1</DistributionRelease>
//   <Architecture>x86_64</Architecture>
//   <InstalledSize>293315967</InstalledSize>
//   <PackageSize>58887876</PackageSize>
//   <PackageHash>b4131d377290e2b6d980db596c77ccc428ee6087</PackageHash>
//   <PackageURI>g/golang/golang-1.7.5-13-1-x86_64.eopkg</PackageURI>
//   <DeltaPackages>
//     <Delta releaseFrom="12">
//       <PackageURI>g/golang/golang-12-13-1-x86_64.delta.eopkg</PackageURI>
//       <PackageSize>47796336</PackageSize>
//       <PackageHash>423f322bf268787466189108b1142222ebfcbfcf</PackageHash>
//     </Delta>
//     <Delta releaseFrom="11">
//       <PackageURI>g/golang/golang-11-13-1-x86_64.delta.eopkg</PackageURI>
//       <PackageSize>47796336</PackageSize>
//       <PackageHash>6c95b53fe40225f072c7400f2774a6898cf8dec1</PackageHash>
//     </Delta>
//     <Delta releaseFrom="10">
//       <PackageURI>g/golang/golang-10-13-1-x86_64.delta.eopkg</PackageURI>
//       <PackageSize>47834176</PackageSize>
//       <PackageHash>c550f5905b540c7e420add3042bf8a1122b3dacf</PackageHash>
//     </Delta>
//     <Delta releaseFrom="9">
//       <PackageURI>g/golang/golang-9-13-1-x86_64.delta.eopkg</PackageURI>
//       <PackageSize>47943656</PackageSize>
//       <PackageHash>75ed3caa07d6b5a11630d1f68377632619cd173b</PackageHash>
//     </Delta>
//   </DeltaPackages>
//   <PackageFormat>1.2</PackageFormat>
//   <Source>
//     <Name>golang</Name>
//     <Homepage>http://golang.org</Homepage>
//     <Packager>
//       <Name>Joshua Strobl</Name>
//       <Email>joshua@stroblindustries.com</Email>
//     </Packager>
//   </Source>
// </Package>

type PISI struct {
	Packages []Package `xml:"Package"`
}

type Package struct {
	Name                string
	Summary             string
	Description         string
	PartOf              string
	License             string
	RuntimeDependencies []Dependency
	Updates             []Update `xml:"History>Update"`
	Source              Source
}

type Dependency struct {
	Name        string `xml:"Dependency"`
	Release     int
	ReleaseFrom int
}

type Update struct {
	Release int `xml:"release,attr"`
	Date    string
	Version string
	Comment string
	Name    string
	Email   string
}

type Source struct {
	Name     string
	Packager Packager
}

type Packager struct {
	Name  string
	Email string
}

func Decode(reader io.Reader) *PISI {
	decoder := xml.NewDecoder(reader)
	pkg := &PISI{}
	decoder.Decode(pkg)

	b, _ := json.MarshalIndent(pkg, "", "  ")
	buff := bytes.NewBuffer(b)
	buff.WriteTo(os.Stdout)
	os.Stdout.WriteString("\n")
	return pkg
}
