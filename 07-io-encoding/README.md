# 07 — I/O, encoding, filesystems

`io.Reader` and `io.Writer` are two methods that the entire Go ecosystem agrees
on. Once you can implement them correctly — including the awkward cases — files,
sockets, compressors, hashers and HTTP bodies all become the same thing.

| # | Exercise | Difficulty | What it drills |
|---|----------|-----------|----------------|
| 1 | `readers` | ★★★ | The full Reader/Writer contract, short reads, `n>0` with `io.EOF` |
| 2 | `framing` | ★★★ | A binary wire format, `io.ReadFull`, checksums, varints, fuzzing |
| 3 | `jsonx` | ★★★ | Custom marshalers, polymorphic decode, streaming `json.Decoder` |
| 4 | `fsx` | ★★★ | `io/fs`, `fstest.MapFS`, atomic writes, reading a file backwards |

## Ideas worth internalising here

**`Read` is allowed to be lazy.** It may return fewer bytes than you asked for,
with a nil error, and it may return data *and* `io.EOF` in the same call. Any
loop of the form `for { n, err := r.Read(p); if err != nil { break }; use(p[:n]) }`
is buggy — it drops the final chunk. Handle `n` first, then `err`.

**`io.ReadFull` and `io.Copy` exist so you don't have to get that right twice.**
`ReadFull` turns "not enough bytes" into `io.ErrUnexpectedEOF`, which is exactly
the distinction a protocol decoder needs.

**Never trust a length prefix.** A four-byte length field from the network can
say 4 GB. Check it against a maximum *before* you allocate. This is a real
denial-of-service class, and the fuzz target in `framing` exists to prove your
decoder can't be crashed by hostile bytes.

**`MarshalJSON` on the value receiver, `UnmarshalJSON` on the pointer.** Get it
backwards and your marshaler is silently skipped when the value isn't
addressable — no error, just wrong output.

**Stream when the input can be large.** `json.Decoder` with `Token`/`More`
processes a 10 GB array in constant memory. `json.Unmarshal` does not.

**Write files atomically.** Temp file in the same directory, `Sync`, `Rename`.
Without the `Sync`, a crash can leave a renamed-but-empty file — the rename is
atomic, the data reaching the disk is not.

**Take `fs.FS`, not a path.** Then `fstest.MapFS` gives you a whole filesystem
in a map literal and your tests need no disk at all.

## Stretch goals

- Add `gzip` to `framing`: compress payloads over a threshold and set a flag bit
  in the type byte. Make sure `ReadFrame` still validates the checksum of the
  *uncompressed* payload.
- Implement `io.ReaderAt` and `io.Seeker` on an in-memory buffer, and use them
  to make `TailLines` work on any `io.ReaderAt` instead of a path.
- Write a `jsonx.Diff(a, b []byte) ([]string, error)` that reports the dotted
  paths where two documents differ, reusing `Flatten`.
