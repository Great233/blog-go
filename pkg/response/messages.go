package response

const (
	Ok            = "ok"
	Fail          = "fail"
	InvalidParams = "请求参数错误"

	UserIsNotExist         = "用户不存在"
	PasswordVerifyFailed = "密码错误"
	TokenGenerateFailed  = "Token生成失败"
	TokenHasExpired      = "Token已过期"
	TokenAuthFailed      = "Token鉴权失败"

	TagIsNotExist       = "标签不存在"
	TagIsAlreadyExist = "标签不存在"
	AddTagFailed      = "添加标签失败"
	EditTagFailed     = "编辑标签失败"
	DeleteTagFailed   = "删除标签失败"

	ArticleIsNotExist           = "文章不存在"
	PathOrTitleIsAlreadyExist = "标题或路径已存在"
	AddArticleFailed          = "添加文章失败"
	EditArticleFailed         = "编辑文章失败"
	DeleteArticleFailed       = "删除文章失败"
)
