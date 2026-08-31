{{template "base" .}}
{{define "content"}}
    <tr>
        <td class="wrapper">
            <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                <tr>
                    <td>
                        <p>👋&nbsp; 你好~ {{.UserName}} ~ </p>
                        <p>💡&nbsp; 感谢您注册账户，很高兴您可以加入我们的大家庭！请激活您的账户。</p>
                        <p>📬&nbsp; 这只需简单一步：点击以下按钮验证您的电子邮件地址。</p>
                        <table role="presentation" border="0" cellpadding="0" cellspacing="0" class="btn btn-primary">
                            <tbody>
                            <tr>
                                <td align="center">
                                    <table role="presentation" border="0" cellpadding="0" cellspacing="0">
                                        <tbody>
                                        <tr>
                                            <td><a href="{{.URL}}" target="_blank">激活账户</a></td>
                                        </tr>
                                        </tbody>
                                    </table>
                                </td>
                            </tr>
                            </tbody>
                        </table>
                        <p>🕹&nbsp; 激活后，您将完成访问！</p>
                        <p>💃&nbsp; 按钮没反应？尝试将此 URL 粘贴到您的浏览器中：<a class='long-url'>{{.URL}}</a></p>
                        <p>😉&nbsp; 我们期待着您的到来！</p>
                    </td>
                </tr>
            </table>
        </td>
    </tr>
{{end}}