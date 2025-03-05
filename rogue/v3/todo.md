# TODO

- [ ] refactor MergeOp
- [ ] add a visLineWeight and totLineWeight to rope for quickly navigating lines
- [ ] revisit caching in github action to speed up again
- [ ] when we implement nested code blocks and quotes in lists should we just add a tab
      value to each type line, code and quote and use that to tab them in even without lists?
- [ ] consistent use of *beforeID and *afterID as well as startID and endID.
      notice that before and after can be nil to indicate the entire document,
      while start and end should never be nil
- [ ] should all the VisLeftOf, TotLeftOf, etc. functions return \*ID instead
      of using ErrorNoLeftTotSibling, etc.?
- [ ] remove start and end id from content addresses
- [ ] remove undoAddress from rewind op, add start and end id to rewind op,
      when applying rewind op, take the current content address to create the redo op.
- [ ] modify selections so that nil means the entire document
- [ ] super formats implmentation. `Format` function accepts format deltas as it currently does, ops are
      now emitted with the current full format(s) as calculated by the nosid. Interval tree
      will now directly store incoming format ops which are no longer deltas. Interval tree
      is modified to include max_seq_id and min_seq_id (for its subtree) for each node. This
      allows us to do a much more efficient search for the max overlapping span less than a given
      id. New core search function is `func MaxOverlappingSpanLessThan(id, startID, endID ID) (formatOp, error)`.
- [ ] factorize types and messages out into submodule so other submodules can use them
      without circular dependencies in the top level rogue package (comes up with formats a lot)
- [ ] factorize formats and abstract its use of rope.GetIndex so that we can test
      it better and migrate to rogue prime easier
- [ ] can make line formats a proper crdt, drop nos lines and just use the formats implementation
- [ ] new nonStickFormatOp and lineFormatOp ops
- [ ] convert deprecated ops in the ToOps function during snapshotting
- [ ] Rogue Prime implementation
- [ ] EnclosingSpanIDV2 performance improvements
- [ ] Make undo/redo respect the current selection
- [ ] replace adjustedParentId function with a ParentID() method on the
      FugueNode that properly adjusts the parentID for the current node
- [ ] replace any use of gods with our own interval tree implementation

- [x] Bug in not sticky formats where they can become sticky
- [x] overhaul interval tree to use printf("%s%s", rev(id), index(startID)) as key
- [x] fix no newline at end of plainline bug
- [x] flatten multiops with a single op in to that op
- [x] replace UndoFormatOp with just formats
- [x] Tease out Op indexing from insert/format/delete so we can bundle many ops
      into a single MultiOp
- [x] Format rewind implementation
- [x] Changing a bullet or number format should update all equally indented
      siblings in the list to the same bullet or number type
- [x] The formats interval tree could now hold the full format derived
      from the nosid instead of just holding the deltas. This is great because
      it'll make undoing formats much easier since the previous format will
      hold the entire state, not just the delta
- [x] SearchOverlapping formats should return an iterator that merges all
      the overlapping spans in lamport sorted order
- [x] Swap out trees instead of slices in the formats interval tree nodes
- [x] Make undo/redo global for all authors and not just the current author
- [x] swap in EnclosingSpanID in RogueEditor
- [x] GetHtml should put ID on inner most tag
- [x] Bug to check current line format when hitting return and decide to
      keep or not
- [x] Bug to check current line format when deleting a newline and deciding
      what format to cleanup
- [x] Bug need to be able to toggle format at end of newline
- [x] Infinite loop bug highlighting across too many format spans in GetCurSpanFormat
- [x] fix remove line format, should just pass empty {} over the wire I think
- [x] fix off by one bug for sticky in getCurrentSpanFormat
- [x] fix off by two bug in toggleCurSpanFormat
- [x] nosid merge function needs to actually merge neighboring spans
      when they share equal formats
- [x] When a line format is passed to the Format function, it should scan
      right to find the correct \n to apply it to. This will make it easier to
      implement line formats in the javascript. What should happen if this
      is done over a selection? It should find all the contained newlines as well
      as the last one to the right of the selection and apply the format to all
      of the lines
- [x] Might want Format() to split sticky/non-sticky formats into two ops
- [x] GetCurSpanFormat(StartID, EndID) (FormatV2, error) to get the current
      format of a selection
- [x] GetCurLineFormat(StartID, EndID)
- [x] nosid should be defined on the rogue and kept up to date with
      calls to Format and formatOp
- [x] GetHtml() should be updated to use the nosid structure
- [x] need a function with nosid to clip it to a StartID and EndID
- [x] figure out how to handle sticky and non-sticky formats in nosid.
      Update: we're splitting out a sticky, not sticky and line nosid
- [x] Swap out FormatV3 for FormatV2 in the rogue
