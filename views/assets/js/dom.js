const loginTerminal = `
<div id="login" class="terminal-box">
                <div class="terminal-head">&lt;LOGIN/&gt;</div>
                <div class="terminal-body">
                    <form id="login-form">
                        <p id="login-text" class="infos-txt">
                            <!-- Nickname row -->
                            <span id="login-intro">
                                <span class="txt command"></span><br>
                                <span class="txt"></span><br>
                                <span class="txt"></span><br>
                                <span class="txt"></span><br>
                                <span class="txt warning"></span><br>
                                <br><br>
                                <span class="txt command"></span><br>
                            </span>

                            <!-- Nickname row -->
                            <span class="auth-row" id="login-username-row">
                                <span class="txt"></span>
                                <span class="txt bold-fluo"></span>
                                <span class="txt"></span>
                                <input type="text" name="login-username" id="login-username" class="auth-input"
                                    autofocus>
                                <span class="txt"></span>
                                <span class="measurer"></span>
                            </span>
                            <br>

                            <!-- Password row -->
                            <span class="auth-row" id="login-password-row">
                                <span class="txt"></span>
                                <span class="txt bold-fluo"></span>
                                <span class="txt"></span>
                                <input type="password" name="login-password" id="login-password" class="auth-input">
                                <span class="txt"></span>
                                <span class="measurer"></span>
                            </span>
                            <br>
                        </p>
                    </form>
                </div>
            </div>
`

const registerTerminal = `
<div id="register" class="terminal-box not-displayed">
                <div class="terminal-head">&lt;REGISTER/&gt;</div>
                <div class="terminal-body">
                    <form id="register-form">
                        <p id="register-text" class="infos-txt">
                            <!-- Nickname row -->
                            <span id="register-intro">
                                <span class="txt command"></span><br>
                                <span class="txt"></span><br>
                                <span class="txt"></span><br>
                                <span class="txt"></span><br>
                                <span class="txt warning"></span><br>
                                <br><br>
                                <span class="txt command"></span><br>
                            </span>

                            <!-- Nickname row -->
                            <span class="auth-row" id="register-username-row">
                                <span class="txt"></span>
                                <span class="txt bold-fluo"></span>
                                <span class="txt"></span>
                                <input type="text" name="register-username" id="register-username" class="auth-input"
                                    autofocus>
                                <span class="txt"></span>
                                <span class="measurer"></span>
                            </span>
                            <br>

                            <!-- Password row -->
                            <span class="auth-row" id="register-password-row">
                                <span class="txt"></span>
                                <span class="txt bold-fluo"></span>
                                <span class="txt"></span>
                                <input type="password" name="register-password" id="register-password" class="auth-input">
                                <span class="txt"></span>
                                <span class="measurer"></span>
                            </span>
                            <br>
                        </p>
                    </form>
                </div>
            </div>
`

var loginView = `
        <div id="auth-left" class="auth-block">
            <div id="welcome" class="terminal-box">
                <div class="terminal-head">&lt;WELCOME/&gt;</div>
                <div class="terminal-body">
                    <p id="welcome-text">
                        <span class="txt"></span><span class="txt bold-fluo"></span><span class="txt"></span><br>
                        <span class="txt"></span><br>
                        <span class="txt bold-fluo"></span><span class="txt"></span>
                    </p>
                </div>
            </div>

            <div id="logs" class="terminal-box">
                <div class="terminal-head">&lt;LOGS/&gt;</div>
                <div class="terminal-body">
                    <p id="logs-text" class="infos-txt">
                        <span class="txt"></span><br>
                    </p>
                </div>
            </div>
        </div>

        <div id="auth-right" class="auth-block">
            <div id="terminal-zone">
                ${loginTerminal}
            </div>
            
            <div id="auth-btns">
                <span id="login-btn" class="auth-btn active"><span class="btn-bg"></span>LOGIN</span>
                <span id="register-btn" class="auth-btn"><span class="btn-bg"></span>REGISTER</span>
            </div>
        </div>
`