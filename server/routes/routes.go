package routes

import (
	"bbsgo/handlers"
	"bbsgo/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.CORS)

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/register", handlers.RegisterWithCode).Methods("POST")
	api.HandleFunc("/send-code", handlers.SendVerificationCode).Methods("POST")
	api.HandleFunc("/login", handlers.Login).Methods("POST")
	api.HandleFunc("/forums", handlers.GetForums).Methods("GET")
	api.HandleFunc("/config", handlers.GetSiteConfig).Methods("GET")
	api.HandleFunc("/topics", handlers.GetTopics).Methods("GET")
	api.HandleFunc("/topics/{id}", handlers.GetTopic).Methods("GET")
	api.HandleFunc("/topics/{id}/posts", handlers.GetPosts).Methods("GET")
	api.HandleFunc("/tags", handlers.GetTags).Methods("GET")
	api.HandleFunc("/tags/search", handlers.SearchTags).Methods("GET")
	api.HandleFunc("/tags/{id}", handlers.GetTag).Methods("GET")
	api.HandleFunc("/announcements", handlers.GetAnnouncements).Methods("GET")
	api.HandleFunc("/users/credit", handlers.GetCreditUsers).Methods("GET")
	api.HandleFunc("/search", handlers.Search).Methods("GET")

	auth := api.PathPrefix("").Subrouter()
	auth.Use(middleware.Auth)

	auth.HandleFunc("/user/profile", handlers.GetProfile).Methods("GET")
	auth.HandleFunc("/user/profile", handlers.UpdateProfile).Methods("PUT")
	auth.HandleFunc("/user/topics", handlers.GetUserTopics).Methods("GET")
	auth.HandleFunc("/user/signin", handlers.SignIn).Methods("POST")
	auth.HandleFunc("/user/signin/status", handlers.GetSignInStatus).Methods("GET")
	auth.HandleFunc("/user/favorites", handlers.GetFavorites).Methods("GET")
	auth.HandleFunc("/user/follows", handlers.GetFollows).Methods("GET")
	auth.HandleFunc("/user/followers", handlers.GetFollowers).Methods("GET")
	auth.HandleFunc("/user/badges", handlers.GetUserBadges).Methods("GET")
	auth.HandleFunc("/user/reports", handlers.GetUserReports).Methods("GET")

	auth.HandleFunc("/topics", handlers.CreateTopic).Methods("POST")
	auth.HandleFunc("/topics/{id}", handlers.UpdateTopic).Methods("PUT")
	auth.HandleFunc("/topics/{id}", handlers.DeleteTopic).Methods("DELETE")

	auth.HandleFunc("/topics/{id}/posts", handlers.CreatePost).Methods("POST")
	auth.HandleFunc("/posts/{id}", handlers.UpdatePost).Methods("PUT")
	auth.HandleFunc("/posts/{id}", handlers.DeletePost).Methods("DELETE")

	auth.HandleFunc("/likes", handlers.CreateLike).Methods("POST")
	auth.HandleFunc("/likes", handlers.DeleteLike).Methods("DELETE")

	auth.HandleFunc("/favorites", handlers.CreateFavorite).Methods("POST")
	auth.HandleFunc("/favorites", handlers.DeleteFavorite).Methods("DELETE")

	auth.HandleFunc("/follows", handlers.CreateFollow).Methods("POST")
	auth.HandleFunc("/follows", handlers.DeleteFollow).Methods("DELETE")
	auth.HandleFunc("/follows/check", handlers.CheckFollow).Methods("GET")

	auth.HandleFunc("/messages", handlers.GetMessages).Methods("GET")
	auth.HandleFunc("/messages", handlers.SendMessage).Methods("POST")
	auth.HandleFunc("/messages/unread-count", handlers.GetUnreadMessageCount).Methods("GET")
	auth.HandleFunc("/messages/read", handlers.MarkMessagesRead).Methods("PUT")
	auth.HandleFunc("/messages/with/{user_id}", handlers.GetMessageConversation).Methods("GET")

	auth.HandleFunc("/notifications", handlers.GetNotifications).Methods("GET")
	auth.HandleFunc("/notifications/unread-count", handlers.GetUnreadNotificationCount).Methods("GET")
	auth.HandleFunc("/notifications/read-all", handlers.MarkAllNotificationsRead).Methods("PUT")

	auth.HandleFunc("/drafts", handlers.GetDrafts).Methods("GET")
	auth.HandleFunc("/drafts", handlers.CreateDraft).Methods("POST")
	auth.HandleFunc("/drafts/{id}", handlers.GetDraft).Methods("GET")
	auth.HandleFunc("/drafts/{id}", handlers.UpdateDraft).Methods("PUT")
	auth.HandleFunc("/drafts/{id}", handlers.DeleteDraft).Methods("DELETE")

	auth.HandleFunc("/reports", handlers.CreateReport).Methods("POST")

	auth.HandleFunc("/badges", handlers.GetBadges).Methods("GET")

	auth.HandleFunc("/upload", handlers.UploadFile).Methods("POST")

	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.Auth)
	admin.Use(middleware.AdminAuth)

	admin.HandleFunc("/users", handlers.GetAdminUsers).Methods("GET")
	admin.HandleFunc("/users/{id}/role", handlers.UpdateUserRole).Methods("PUT")
	admin.HandleFunc("/users/{id}/ban", handlers.BanUser).Methods("PUT")
	admin.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

	admin.HandleFunc("/forums", handlers.CreateForum).Methods("POST")
	admin.HandleFunc("/forums/{id}", handlers.UpdateForum).Methods("PUT")
	admin.HandleFunc("/forums/{id}", handlers.DeleteForum).Methods("DELETE")

	admin.HandleFunc("/tags", handlers.GetAdminTags).Methods("GET")
	admin.HandleFunc("/tags", handlers.CreateTag).Methods("POST")
	admin.HandleFunc("/tags/merge", handlers.MergeTags).Methods("POST")
	admin.HandleFunc("/tags/{id}", handlers.UpdateTag).Methods("PUT")
	admin.HandleFunc("/tags/{id}", handlers.DeleteTag).Methods("DELETE")

	admin.HandleFunc("/topics", handlers.GetAdminTopics).Methods("GET")
	admin.HandleFunc("/topics/{id}", handlers.DeleteAdminTopic).Methods("DELETE")

	admin.HandleFunc("/posts", handlers.GetAdminPosts).Methods("GET")
	admin.HandleFunc("/posts/{id}", handlers.DeleteAdminPost).Methods("DELETE")

	admin.HandleFunc("/reports", handlers.GetAdminReports).Methods("GET")
	admin.HandleFunc("/reports/{id}/handle", handlers.HandleReport).Methods("PUT")

	admin.HandleFunc("/announcements", handlers.CreateAnnouncement).Methods("POST")
	admin.HandleFunc("/announcements/{id}", handlers.UpdateAnnouncement).Methods("PUT")
	admin.HandleFunc("/announcements/{id}", handlers.DeleteAnnouncement).Methods("DELETE")

	admin.HandleFunc("/config", handlers.UpdateSiteConfig).Methods("PUT")
	admin.HandleFunc("/forum-categories", handlers.GetAllForumCategories).Methods("GET")
	admin.HandleFunc("/forum-categories", handlers.CreateForumCategory).Methods("POST")
	admin.HandleFunc("/forum-categories/{id}", handlers.UpdateForumCategory).Methods("PUT")
	admin.HandleFunc("/forum-categories/{id}", handlers.DeleteForumCategory).Methods("DELETE")
	admin.HandleFunc("/change-password", handlers.ChangeAdminPassword).Methods("POST")

	return r
}
