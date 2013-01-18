package psi

type Descriptor []byte

func (d Descriptor) Tag() byte {
	return d[0]
}

func (d Descriptor) Data() []byte {
	return d[2 : 2+d[1]]
}

type DescriptorList []byte

// Pop returns first descriptor in d and remaining descriptors in rdl. If an
// error occurs it returns d == nil. If there is no more descriptors
// len(rdl) == 0
func (dl DescriptorList) Pop() (d Descriptor, rdl DescriptorList) {
	if len(dl) < 2 {
		return
	}
	l := int(dl[1]) + 2
	if len(dl) < l {
		return
	}
	d = Descriptor(dl[2:l])
	rdl = dl[l:]
	return
}

// Append adds d to the end of the dl. It works like Go append function so need
// to be used in this way:
//     dl = dl.Append(d)
func (dl DescriptorList) Append(d Descriptor) DescriptorList {
	return append(dl, d...)
}
