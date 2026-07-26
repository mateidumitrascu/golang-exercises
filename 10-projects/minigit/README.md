# Project: mini-git

**Modules it exercises:** 02 (bytes), 04 (errors), 07 (I/O, fs), 09 (algorithms).
**Rough size:** 700–1000 lines. **Difficulty:** ★★★

A content-addressable version control system. Not a git clone — the *idea* of
git, implemented from scratch, in a way you will never forget.

## The spec

```
minigit init
minigit add <path>...
minigit commit -m "message"
minigit log
minigit status
minigit diff [<commit>]
minigit checkout <commit>
minigit cat-file <hash>
```

Storage, under `.minigit/`:

```
objects/ab/cdef123...     zlib-compressed objects, named by SHA-256 of the content
refs/heads/main           a file containing a commit hash
HEAD                      "ref: refs/heads/main"
index                     the staging area
```

Three object types, each stored as `"<type> <length>\0<payload>"`:

- **blob** — a file's bytes.
- **tree** — a sorted directory listing: mode, name, and the hash of a blob or
  another tree.
- **commit** — a tree hash, parent hash(es), author, timestamp, message.

## Requirements

1. **Content addressing.** The hash of an object is the hash of its bytes, so
   identical files stored twice occupy one object. Prove it in a test.
2. **Immutability.** Objects are never modified, only added. Write them
   atomically (temp file + rename — module 07).
3. **Trees are recursive.** `commit` builds a tree of trees from the index.
4. **`log` walks the parent chain** and prints hash, date, message.
5. **`status`** compares the working directory, the index and HEAD's tree, and
   reports the three states (staged / modified / untracked).
6. **`diff`** produces a unified diff. The algorithm is the LCS from
   `02-strings-bytes/edit` — you have already written the hard part.
7. **`checkout`** restores a commit's tree into the working directory and
   refuses to run with uncommitted changes.

## Milestones

1. `hash-object` and `cat-file`: write a blob, read it back. This is the whole
   idea in 50 lines.
2. The index (a simple sorted `path -> hash` file) and `add`.
3. Trees and `commit`. Draw the object graph on paper first.
4. `log`, then `status`.
5. `diff` with the LCS.
6. `checkout` and branch refs.

## The interesting problem

Trees. A commit points to one tree; a tree points to blobs and other trees.
Building that from a flat index means grouping paths by directory prefix and
recursing bottom-up. Once it works, notice that two commits sharing an unchanged
subdirectory *share the same tree object* — that is the whole reason git is fast,
and it falls out of content addressing for free.

## Stretch

- Branches and a three-way merge with conflict markers.
- Packing: store objects as deltas against a base, like git's packfiles.
- `gc`: find unreachable objects from the refs and delete them (a graph
  traversal — module 09).
