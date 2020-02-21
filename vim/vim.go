package vim

import "io/ioutil"

const vimConfig = `
:imap jj <Esc>
set tabstop=4 softtabstop=0 expandtab shiftwidth=2 smarttab

if executable('ag')
	set grepprg=ag\ --nogroup\ --nocolor
endif

" make the quickfix window automatically close when selecting a file with enter key
:autocmd FileType qf nnoremap <buffer> <CR> <CR>:cclose<CR>

" automatically set working directory to match file. this makes grep! work easily
set autochdir

" TODO use vim regex to highlight [[stuff]]
" highlight link CardsLink markdownLinkText

" jump among [[links]]
:nnoremap J /\[\[.*\]\]<cr>:noh<cr>ll

" pop open a list of file matches based on word under cursor
:nnoremap K :silent execute "grep! -R " . shellescape(expand("<cword>")) . " ."<cr>:copen<cr>
`

func CreateConfig() string {
	file, err := ioutil.TempFile("/tmp", "cards-config-*.vimrc")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(vimConfig)
	if err != nil {
		panic(err)
	}
	return file.Name()
}
