package psi

// SectionReplacer is an interface that wraps the ReplaceSection method. After
// ReplaceSection call caller should not reffer to s content any more.
// If ReplaceSection returns an error it is guaranteed that r == s (but content
// reffered by s can be modified). Generally ReplaceSection should be used in
// this way:
//
//    s, err = q.ReplaceSection(s)
//    if err != nil {
//        ...
//    }
type SectionReplacer interface {
	// ReplaceSection consumes secion reffered by s an returns other section
	// reffered by r.
	ReplaceSection(s Section) (r Section, e error)
}
