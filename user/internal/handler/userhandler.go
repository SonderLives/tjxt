package handler

import (
	"net/http"
	"strconv"

	"common/xerr"
	"user/internal/logic"
	"user/internal/svc"
	"user/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func StudentRegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StudentFormDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := ValidateStudentForm(&req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewStudentRegisterLogic(r.Context(), svcCtx)
		resp, err := l.StudentRegister(&req)
		writeResult(w, r, resp, err)
	}
}

func StudentPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StudentFormDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := ValidateStudentForm(&req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewStudentPasswordLogic(r.Context(), svcCtx)
		resp, err := l.StudentPassword(&req)
		writeResult(w, r, resp, err)
	}
}

func StudentPageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserPageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewStudentPageLogic(r.Context(), svcCtx)
		resp, err := l.StudentPage(&req)
		writeResult(w, r, resp, err)
	}
}

func TeacherPageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserPageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewTeacherPageLogic(r.Context(), svcCtx)
		resp, err := l.TeacherPage(&req)
		writeResult(w, r, resp, err)
	}
}

func StaffPageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserPageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewStaffPageLogic(r.Context(), svcCtx)
		resp, err := l.StaffPage(&req)
		writeResult(w, r, resp, err)
	}
}

func UserMeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUserMeLogic(r.Context(), svcCtx)
		resp, err := l.UserMe()
		writeResult(w, r, resp, err)
	}
}

func UserCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := ValidateUserDTO(&req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewUserCreateLogic(r.Context(), svcCtx)
		resp, err := l.UserCreate(&req)
		writeResult(w, r, resp, err)
	}
}

func UserUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFormDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := ValidateUserDTO(&req.UserDTO); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewUserUpdateLogic(r.Context(), svcCtx)
		resp, err := l.UserUpdate(&req)
		writeResult(w, r, resp, err)
	}
}

func UserCheckCellphoneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserCheckCellphoneReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewUserCheckCellphoneLogic(r.Context(), svcCtx)
		resp, err := l.UserCheckCellphone(&req)
		writeResult(w, r, resp, err)
	}
}

func UserGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserIdReq
		id, err := strconv.ParseInt(pathvar.Vars(r)["id"], 10, 64)
		if err != nil || id <= 0 {
			httpx.ErrorCtx(r.Context(), w, xerr.New(xerr.CodeBadRequest, "路径参数 id 非法"))
			return
		}
		req.Id = id
		l := logic.NewUserGetLogic(r.Context(), svcCtx)
		resp, err := l.UserGet(&req)
		writeResult(w, r, resp, err)
	}
}

func UserAdminUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 先解析 path param，再解析 body（避免 httpx.Parse 消费 body）
		id, err := strconv.ParseInt(pathvar.Vars(r)["id"], 10, 64)
		if err != nil || id <= 0 {
			httpx.ErrorCtx(r.Context(), w, xerr.New(xerr.CodeBadRequest, "路径参数 id 非法"))
			return
		}
		var req types.UserIdReq
		req.Id = id

		var body types.UserDTO
		if err := httpx.ParseJsonBody(r, &body); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := ValidateUserDTO(&body); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserAdminUpdateLogic(r.Context(), svcCtx)
		resp, err := l.UserAdminUpdate(&req, &body)
		writeResult(w, r, resp, err)
	}
}

func UserResetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserIdReq
		id, err := strconv.ParseInt(pathvar.Vars(r)["id"], 10, 64)
		if err != nil || id <= 0 {
			httpx.ErrorCtx(r.Context(), w, xerr.New(xerr.CodeBadRequest, "路径参数 id 非法"))
			return
		}
		req.Id = id
		l := logic.NewUserResetPasswordLogic(r.Context(), svcCtx)
		resp, err := l.UserResetPassword(&req)
		writeResult(w, r, resp, err)
	}
}

func UserStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if req.Id <= 0 || req.Status <= 0 {
			httpx.ErrorCtx(r.Context(), w, xerr.New(xerr.CodeBadRequest, "路径参数 id/status 必须大于 0"))
			return
		}
		l := logic.NewUserStatusLogic(r.Context(), svcCtx)
		resp, err := l.UserStatus(&req)
		writeResult(w, r, resp, err)
	}
}