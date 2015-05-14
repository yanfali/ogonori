package oschema

import "fmt"

//
// This file holds LINK type datastructures.
// Namely, for LINK, LINKLIST (LINKSET) and LINKMAP.
//

//
// DOCUMENT ME
//
type OLink struct {
	RID    ORID
	Record *ODocument
}

func (lnk *OLink) String() string {
	recStr := "<nil>"
	if lnk.Record != nil {
		// fields are not shown to avoid infinite loops when there are circular links
		recStr = lnk.Record.StringNoFields()
	}
	return fmt.Sprintf("<OLink RID: %s, Record: %s>", lnk.RID, recStr)
}

// ------

//
// OLinkBag can have a tree-based or an embedded representation.
//
// Embedded stores its content directly to the document that owns it.
// It is used when only small numbers of links are stored in the bag.
//
// The tree-based implementation stores its content in a separate data
// structure called on OSBTreeBonsai on the server. It fits great for cases
// when you have a large number of links.  This is used to efficiently
// manage relationships (particularly in graph databases).
//
// The OLinkBag struct corresponds to ORidBag in Java client codebase.
//
//
type OLinkBag struct {
	Links             []*OLink
	Size              int
	CollectionPointer *treeCollectionPointer
}

// // EXPTAL
// type OLinkBag struct {
// 	Links             []*OLink
// 	Size              int
// 	TreeOLinkBag
// }

// type TreeOLinkBag struct {
// 	fileId     int64
// 	pageIndex  int64
// 	pageOffset int32
// }
// // END EXPTAL

type treeCollectionPointer struct {
	fileId     int64
	pageIndex  int64
	pageOffset int32
}

func (lb OLinkBag) GetFileID() int64 {
	return lb.CollectionPointer.fileId
}

func (lb OLinkBag) GetPageIndex() int64 {
	return lb.CollectionPointer.pageIndex
}

func (lb OLinkBag) GetPageOffset() int32 {
	return lb.CollectionPointer.pageOffset
}

//
// NewOLinkBag constructor is called with all the OLink
// objects precreated. Usually appropriate when dealing
// with an embedded LinkBag.
//
func NewOLinkBag(links []*OLink) *OLinkBag {
	return &OLinkBag{Links: links, Size: len(links)}
}

//
// NewTreeOLinkBag constructor is called for remote tree-based
// LinkBags.  This is called by the Deserializer when all it knows
// is the pointer reference to the LinkBag on the remote server.
//
// The OLinkBag returned does not yet know the size of the LinkBag
// nor know what the OLinks are.
//
func NewTreeOLinkBag(fileId int64, pageIdx int64, pageOffset int32, size int32) *OLinkBag {
	treeptr := treeCollectionPointer{
		fileId:     fileId,
		pageIndex:  pageIdx,
		pageOffset: pageOffset,
	}

	return &OLinkBag{CollectionPointer: &treeptr, Size: int(size)}
}

// // ------

// type OLinkMap struct {
// 	Links map[string]*OLink
// }