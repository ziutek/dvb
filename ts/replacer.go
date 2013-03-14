package ts

// PktReplacer is an interface that wraps the ReplacePkt method. After
// ReplacePkt call caller should not reffer to p content any more.
// If ReplacePkt returns an error it is guaranteed that r == p (but content
// reffered by p can be modified). Generally ReplacePkt should be used in this
// way:
//
//    p, err = q.ReplacePkt(p)
//    if err != nil {
//        ...
//    }
type PktReplacer interface {
	// ReplacePkt consumes packet reffered by p an returns other packet reffered
	// by r.
	// If you use ReplacePkt for reading ErrSync or dvb.ErrOverflow are not
	// fatal errors. You can still call ReplacePkt after obtaining such errors.
	// If you use ReplacePkt for writing, any error is probably a problem.
	ReplacePkt(p *ArrayPkt) (r *ArrayPkt, e error)
}

// PktReaderAsReplacer converts any PktReader to PktReplacer
type PktReaderAsReplacer struct {
	R PktReader
}

func (r PktReaderAsReplacer) ReplacePkt(p *ArrayPkt) (*ArrayPkt, error) {
	err := r.R.ReadPkt(p)
	return p, err
}

// PktWriterAsReplacer converts any PktWriter to PktReplacer
type PktWriterAsReplacer struct {
	W PktWriter
}

func (r PktWriterAsReplacer) ReplacePkt(p *ArrayPkt) (*ArrayPkt, error) {
	err := r.W.WritePkt(p)
	return p, err
}
