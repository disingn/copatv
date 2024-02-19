package main

import (
	"changeme/sever"
	"context"
	"fmt"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (a *App) GetVersion(auth, tool, name string) string {
	err, d := sever.Getauth(name, auth)
	if err != nil {
		return fmt.Sprintf("授权失败,失败原因:%s", err)
	}
	if tool == "jetbrains" {
		err := sever.SetJbHost(d.Token, d.Url, name)
		if err != nil {
			return fmt.Sprintf("授权失败,失败原因:%s", err)
		}
	} else if tool == "vscode" {
		e := fmt.Sprintf("https://%s:%s@%s", name, d.Token, strings.TrimPrefix(d.Url, "https://"))
		err := sever.SetSetting(d.Url, e)
		if err != nil {
			return fmt.Sprintf("授权失败,失败原因:%s", err)
		}
	}
	return fmt.Sprintf("%s,授权成功,过期时间为%s", tool, d.ExpireTime)
}

func (a *App) OpenUrl(url string) {
	_ = sever.OpenBrowser(url)

}
