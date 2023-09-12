function copy2Clipboard(content) {
    if (window.navigator && window.navigator.clipboard) {
        navigator.clipboard.writeText(content).then(() => {
            console.log('success')
        }, () => {
            console.log('failed')
            console.log(arguments)
        }).catch(e => {
            console.log(e)
        })
    } else {
        var i = document.createElement("input")
        i.setAttribute("style", "position:absolute;z-index:-1;width:1px;height:1px;top:-1px;opacity:0;-webkit-user-select: text;")
        document.body.appendChild(i)
        i.value = content
        var t = i.contentEditable,
            n = i.readOnly,
            o = document.createRange();
        i.contentEditable = !0
        i.readOnly = !1
        o.selectNodeContents(i);
        var a = window.getSelection();
        a.removeAllRanges()
        a.addRange(o)
        i.setSelectionRange(0, 999999)
        i.contentEditable = t
        i.readOnly = n
        window.document.execCommand("Copy")
        i.blur()
        document.body.removeChild(i)
        i = null
    }
}
function scriptLoader(url) {
    return new Promise((resolve, reject) => {
        const script = document.createElement('script')
        script.type = 'text/javascript'
        script.async = true
        script.src = url
        if (script.readyState) {
            script.onreadystatechange = () => {
                if (script.readyState === 'loaded' || script.readyState === 'complete') {
                    script.onreadystatechange = null
                    resolve()
                }
            }
        } else {
            script.onload = resolve
        }
        document.body.appendChild(script)
    })
}
function loadJsonToVars(jsonUrl) {
    return new Promise((resolve, reject) => {
        if (!jsonUrl) reject;
        fetch(jsonUrl).then(r => r.text()).then(jsonText => {
            let variable = JSON.parse(jsonText);
            resolve(variable)
        }).catch(e => {
            reject()
        })
    })
}
function dispatchCustomEvent(eventName, payload) {
    const evt = new CustomEvent(eventName, { detail: payload });
    document.dispatchEvent(evt);
}
document.addEventListener("alpine:init", () => {
    // Global data
    Alpine.data('root', () => ({
        loading: false,
        last_config: "",
        changed_config: false,
        saving: false,
        schemaUrl: "",
        completionProvide: null,
        __initMonacoEditor() {
            this.loading = true;
            Promise.all([
                scriptLoader("https://cdn.jsdelivr.net/npm/monaco-editor@0.40.0/min/vs/loader.js"),
                scriptLoader("https://cdn.jsdelivr.net/npm/monaco-editor@0.40.0/min/vs/editor/editor.main.nls.js"), // FIXME: Uncaught ReferenceError: define is not defined at editor.main.nls.js:11:1
                scriptLoader("https://cdn.jsdelivr.net/npm/monaco-editor@0.40.0/min/vs/editor/editor.main.js"),
            ]).then((res) => {
                window.onload = () => {
                    this.loading = false;
                    monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
                        validate: true,
                        enableSchemaRequest: true
                    })
                    this.completionProvide = monaco.languages.registerCompletionItemProvider('json', {
                        provideCompletionItems: (model, position) => {
                            // find out if we are completing a property in the 'dependencies' object.
                            var textUntilPosition = model.getValueInRange({
                                startLineNumber: 1,
                                startColumn: 1,
                                endLineNumber: position.lineNumber,
                                endColumn: position.column,
                            });
                            var appsHttpServersMatch = textUntilPosition.match(
                                /"servers": \{?/
                            );
                            var appsMatch = textUntilPosition.match(
                                /"apps": \{?/
                            );
                            var appsHttpMatch = textUntilPosition.match(
                                /"http": \{?/
                            );
                            var appsHttpHandleMatch = textUntilPosition.match(
                                /"handle": \[?/g
                            );
                            var appsTlsMatch = textUntilPosition.match(
                                /"tls": \{?/
                            );
                            var appsLayer4Match = textUntilPosition.match(
                                /"layer4": \{?/
                            );
                            var appsTlsIssuresMatch = textUntilPosition.match(
                                /"issures": \[?/
                            );
                            var appsHttpMatcherMatch = textUntilPosition.match(
                                /"match": \[?/
                            );

                            if (appsMatch) {
                                if (appsLayer4Match) {
                                    return {
                                        suggestions: [
                                            {
                                                label: 'App Layer4',
                                                kind: monaco.languages.CompletionItemKind.Function,
                                                documentation: "",
                                                insertText: `"srv_name": {
"listen": [
    "0.0.0.0:\${1}"
],
"routes": [
    {
        "handle": [
            {
                "handler": "proxy",
                "upstreams": [
                    {
                        "dial": [
                            "127.0.0.1:\${2}"
                        ]
                    }
                ]
            }
        ]
    }
]
}`,
                                                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                            },
                                        ]
                                    }
                                }
                                if (appsTlsMatch) {
                                    if (appsTlsIssuresMatch) {
                                        return {
                                            suggestions: [
                                                {
                                                    label: 'App TLS Issurer (ZeroSSL + Cloudflare = Wildcard)',
                                                    kind: monaco.languages.CompletionItemKind.Function,
                                                    documentation: "",
                                                    insertText: `{
"ca": "https://acme.zerossl.com/v2/DV90",
"challenges": {
    "dns": {
        "propagation_delay": 60,
        "provider": {
            "api_token": "\${1}",
            "name": "cloudflare"
        },
        "ttl": 60
    }
},
"email": "",
"external_account": {
    "key_id": "\${3}",
    "mac_key": "\${2}"
},
"module": "zerossl"
}`,
                                                    insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                                },
                                            ]
                                        }
                                    }
                                    return {
                                        suggestions: [
                                            // tls configuration
                                            {
                                                label: 'App TLS Issurer ACME',
                                                kind: monaco.languages.CompletionItemKind.Function,
                                                documentation: '`apps.tls',
                                                insertText: `{
"automation": {
    "policies": [
        {
            "issuers": [
                {
                    "email": "\${1}",
                    "module": "acme"
                }
            ],
            "subjects": [\${2}]
        }
    ]
},
"certificates": {
    "automate": [\${0}]
}
}`,
                                                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                            },
                                        ]
                                    }
                                }
                                if (appsHttpMatch) {
                                    if (appsHttpServersMatch) {
                                        if (appsHttpMatcherMatch) {
                                            return {
                                                suggestions: [
                                                    {
                                                        label: 'HTTP Server match rule - remote_ip',
                                                        kind: monaco.languages.CompletionItemKind.Function,
                                                        documentation: "",
                                                        insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                                        insertText: `{"remote_ip": {"ranges": [\${0}]}}`,
                                                    },
                                                    {
                                                        label: 'HTTP Server match rule - host',
                                                        kind: monaco.languages.CompletionItemKind.Function,
                                                        documentation: "",
                                                        insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                                        insertText: `{"host": ["\${0}"]}`,
                                                    }
                                                ]
                                            }
                                        }
                                        if (appsHttpHandleMatch) {
                                            return {
                                                suggestions: [
                                                    {
                                                        label: 'HTTP Server handler (reverse_proxy)',
                                                        kind: monaco.languages.CompletionItemKind.Function,
                                                        documentation: '`apps.http.servers[srv_id]`',
                                                        insertText: `{
"group": "",
"handle": [
    {
        "handler": "reverse_proxy",
        "upstreams": [
            {
                "dial": "127.0.0.1:\${1}"
            }
        ]
    }
],
"match": [
    {"host": ["\${0}"]}
]
},`,
                                                        insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                                    },
                                                    {
                                                        label: 'HTTP Server handler (file_server)',
                                                        kind: monaco.languages.CompletionItemKind.Function,
                                                        documentation: '`apps.http.servers[srv_id]`',
                                                        insertText: [
                                                            '{',
                                                            '\t"handler": "file_server",',
                                                            '\t"hide": ["${1}"],',
                                                            '\t"root": "${0}"',
                                                            '}'].join('\n'),
                                                        insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                                    },
                                                    {
                                                        label: 'HTTP Server handler (webdav)',
                                                        kind: monaco.languages.CompletionItemKind.Function,
                                                        documentation: '`apps.http.servers[srv_id]`',
                                                        insertText: `{
"handler": "webdav",
"root": "\${0}",
"prefix": ""
}`,
                                                        insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                                    },
                                                    {
                                                        label: 'HTTP Server handler (http_basic)',
                                                        kind: monaco.languages.CompletionItemKind.Function,
                                                        documentation: '`apps.http.servers[srv_id]`',
                                                        insertText: `{
"handler": "authentication",
"providers": {
    "http_basic": {
        "accounts": [
            {
                "password": "\${2}",
                "username": "\${1}"
            }
        ],
        "hash": {
            "algorithm": "bcrypt"
        }
    }
}
}`,
                                                        insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                                    },
                                                    {
                                                        label: 'HTTP Server handler (forward_proxy)',
                                                        kind: monaco.languages.CompletionItemKind.Function,
                                                        documentation: '`apps.http.servers[srv_id]`',
                                                        insertText: `{
"group": "naiveproxy",
"handle": [
    {
        "auth_pass_deprecated": "\${2}",
        "auth_user_deprecated": "\${1}",
        "handler": "forward_proxy",
        "hide_ip": true,
        "hide_via": true,
        "probe_resistance": {
            "domain": "\${0}"
        }
    }
]
}`,
                                                        insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                                    },
                                                    {
                                                        label: 'HTTP Server handler (static_response)',
                                                        kind: monaco.languages.CompletionItemKind.Function,
                                                        documentation: '`apps.http.servers[srv_id]`',
                                                        insertText: [
                                                            '{',
                                                            '\t"handler": "static_response",',
                                                            '\t"status_code": 200,',
                                                            '\t"body": "${0}"',
                                                            '}'].join('\n'),
                                                        insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                                    },
                                                    {
                                                        label: 'HTTP Server handler (subroute)',
                                                        kind: monaco.languages.CompletionItemKind.Function,
                                                        documentation: '`apps.http.servers[srv_id]`',
                                                        insertText: [
                                                            '{',
                                                            '\t"handler": "subroute",',
                                                            '\t"routes": [${0}]',
                                                            '}'].join('\n'),
                                                        insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                                    },
                                                ]
                                            }
                                        }

                                        return {
                                            suggestions: [
                                                // http.handlers
                                                {
                                                    label: 'HTTP Server - Wildcard',
                                                    kind: monaco.languages.CompletionItemKind.Function,
                                                    documentation: '`apps.http.servers[srv_id]`',
                                                    insertText: `"srv0": {
"listen": [
    ":\${1}"
],
"protocols": ["h1","h2","h3"],
"routes": [{
    "group": "wildcard match",
    "handle": [{
        "handler": "subroute",
        "routes": [{
            "group": "match example1",
            "handle": [\${0}],
            "match": [{
                "host": ["\${3}"]
            }]
        }]
    }],
    "match": [{
        "host": [
            "\${2}"
        ]
    }]
}]
}`,
                                                    insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                                }
                                            ]
                                        };
                                    }
                                }
                            }


                            return {
                                suggestions: [
                                    {
                                        label: 'Logging Default Configuration',
                                        kind: monaco.languages.CompletionItemKind.Function,
                                        documentation: "",
                                        insertText: `"logging": {
"logs": {
    "default": {
        "level": "INFO",
        "writer": {
            "filename": "/var/log/caddy/access.log",
            "output": "file"
        }
    }
}
}`,
                                        insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                                    },
                                ]
                            };
                        }
                    });
                    let editor = monaco.editor.create(document.getElementById('config-editor'), {
                        value: '{}',
                        language: 'json',
                    })
                    window.CADDY_CONFIG_EDITOR = editor
                    this.__loadServerConfig().then(res => {
                        if (res.ok) {
                            this.last_config = res.data;
                            editor.getModel().setValue(JSON.stringify(JSON.parse(res.data), null, 4))
                            editor.getModel().onDidChangeContent((event) => {
                                this.completionProvide.dispose()
                                this.changed_config = false;
                                try {
                                    let currentJSON = JSON.stringify(editor.getModel().getValue())
                                    if (currentJSON !== this.last_config) this.changed_config = true;
                                } catch (e) {
                                    console.log(e)
                                }
                            });
                        } else {
                            alert(res.msg)
                        }
                    });
                }
            })
        },
        __setMonacoEditorValue(json) {
            console.log(this.editor)
            this.editor.getModel().setValue(json)
        },
        __loadServerConfig() {
            return fetch("/api/load").then(r => r.json())
        },
        __updateServerConfig() {
            if (!this.changed_config) return;
            let currentConfig = window.CADDY_CONFIG_EDITOR.getModel().getValue();
            try {
                let obj = JSON.parse(currentConfig)
                if (obj["$schema"]) delete obj["$schema"]
                currentConfig = JSON.stringify(obj);
            } catch (e) {
                alert(e.message);
                return;
            }
            this.saving = true;
            fetch("/api/save", {
                method: "POST",
                body: currentConfig
            }).then(r => r.json()).then(res => {
                if (res.ok) {
                    this.last_config = JSON.stringify(JSON.parse(currentConfig))
                    this.changed_config = false;
                } else {
                    alert(res.msg)
                }
            }).finally(() => {
                this.saving = false;
            })
        },
        __discardChanges() {
            window.CADDY_CONFIG_EDITOR.getModel().setValue(JSON.stringify(JSON.parse(this.last_config), null, 4))
        },
        init() {
            this.__initMonacoEditor();
            this.schemaUrl = window.location.href + (window.location.href.endsWith("/") ? "caddy_schema.json" : "/caddy_schema.json")
        }
    }))

    // Toolbox - modalBcrypt
    Alpine.data("modalBcrypt", () => ({
        // show: false,
        data: {
            salt: "",
            origin: "",
            hash: ""
        },
        bcrypt: null,
        // toggle() {
        //     this.show = !this.show
        // },
        copy2Clipboard(e) {
            if (this.data.hash) {
                copy2Clipboard(this.data.hash);
                alert("The hashed value has been copied to clipboard.")
            } else {
                alert("No data yet.")
            }
            e.target.blur();
        },
        init() {
            this.bcrypt = dcodeIO.bcrypt;
            this.data.salt = this.bcrypt.genSaltSync(10);
            // document.addEventListener("custom__toolbox__show:modalBcrypt", () => {
            //     this.show = true;
            // })
            this.$watch('data.origin', value => {
                this.data.hash = this.bcrypt.hashSync(value, this.data.salt)
            })
        },
    }))

    Alpine.data('modalStepEditor', () => ({
        // data
        totalSteps: 5,
        currentStep: 1,
        form: {
            overwrite: false,
        }
        // methods

        // lifecycles
    }))
})