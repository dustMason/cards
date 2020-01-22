# cards

My personal CLI tool for keeping notes. Written in go. It depends on these external binaries:

* [fzf](https://github.com/junegunn/fzf)
* [ag](https://github.com/ggreer/the_silver_searcher)

They can be installed with `brew install fzf the_silver_searcher`

All content is stored in `~/.cards`.

---

### TODO

- [ ] `browse` should have a multi select key, with some actions to take on selected items (archive, rm)
- [ ] `browse` should have some single file actions (gist, run, rename)
- [ ] `browse` should have a shortcuts help overlay
- [ ] jump to `search` from `browse` by typing `/`
- [ ] `browse` list pane should include a short preview (first 2 lines?)
- [ ] the tool should manage a git repo for the ~/.cards dir?
