package prism

type ArticleDto struct {
	ArticleType    string            `json:"type"`
	ArticleContent ArticleContentDto `json:"content"`
}

type ArticleContentDto struct {
	LinkURL      string `json:"link_url,omitempty"`
	GistID       string `json:"gist_id,omitempty"`
	MarkdownText string `json:"md_text,omitempty"`
}

func (d *ArticleDto) update(a *Article) {

	switch d.ArticleType {
	case "link":
		a.initLinkArticle(map[string]interface{}{
			"LinkURL": d.ArticleContent.LinkURL,
		})
	case "gist":
		a.initGistArticle(map[string]interface{}{
			"GistID": d.ArticleContent.GistID,
		})

	case "markdown":
		a.initMarkdownArticle(map[string]interface{}{
			"MarkdownText": d.ArticleContent.MarkdownText,
		})
	}

}
