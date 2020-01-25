# cards

My personal CLI tool for keeping notes. Written in go. It depends on these external binaries:

* [fzf](https://github.com/junegunn/fzf)
* [ag](https://github.com/ggreer/the_silver_searcher)

They can be installed with `brew install fzf the_silver_searcher`

All content is stored in `~/.cards`.

---

### TODO

- [ ] `browse` single file actions 
    - [x] pbcopy
    - [x] archive
    - [x] rename
    - [ ] gist
    - [ ] run
- [ ] instead of switching to a new screen, `browse` `rename` should put the input in the table cell
- [ ] jump to `search` from `browse` by typing `/`
- [ ] `new` command should prompt for a file suffix. input field overlay?
- [ ] the tool should manage a git repo for the ~/.cards dir?
- [ ] shortcut key to open markdown file in browser? (render html to tempfile?)
- [ ] a solution for editing notes on mobile and syncing
- [ ] tags and linking!
