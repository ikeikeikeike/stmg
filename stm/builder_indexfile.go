package stm

import "bytes"

// NewBuilderIndexfile returns the created the BuilderIndexfile's pointer
func NewBuilderIndexfile(opts *Options, loc *Location) *BuilderIndexfile {
	return &BuilderIndexfile{opts: opts, loc: loc}
}

// BuilderIndexfile provides implementation for the Builder interface.
type BuilderIndexfile struct {
	opts     *Options
	loc      *Location
	content  []byte
	linkcnt  int
	totalcnt int
	written  bool
}

// Add method joins old bytes with creates bytes by it calls from Sitemap.Finalize method.
func (b *BuilderIndexfile) Add(link interface{}) BuilderError {
	bldr := link.(*BuilderFile)
	bldr.Write()

	smu := NewSitemapIndexURL(b.opts, URL{{"loc", bldr.loc.URL()}})
	b.content = append(b.content, smu.XML()...)

	b.totalcnt += bldr.linkcnt
	b.linkcnt++

	// If not consolidating indexes, write immediately
	if !b.opts.consolidateIndex {
		b.Write()
		b.content = []byte{} // Clear content after writing
		b.linkcnt = 0
	}
	return nil
}

// Content and BuilderFile.Content are almost the same behavior.
func (b *BuilderIndexfile) Content() []byte {
	return b.content
}

// XMLContent and BuilderFile.XMLContent share almost the same behavior.
func (b *BuilderIndexfile) XMLContent() []byte {
	c := bytes.Join(bytes.Fields(IndexXMLHeader), []byte(" "))
	c = append(append(c, b.Content()...), IndexXMLFooter...)

	return c
}

// Write and Builderfile.Write are almost the same behavior.
func (b *BuilderIndexfile) Write() {
	// If consolidating indexes and already written, skip
	if b.opts.consolidateIndex && b.written {
		return
	}

	c := b.XMLContent()
	b.loc.Write(c, b.linkcnt)
	
	if b.opts.consolidateIndex {
		b.written = true
	}
}
