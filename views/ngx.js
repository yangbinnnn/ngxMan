var sites = [
    {title: 'www.cloudhua.com'},
    {title: 'git.cloudhua.com'},
    {title: 'docker.cloudhua.com'},
];

var sitesMap = new Map();
sitesMap.set("www.cloudhua.com", `
                location / {
                    proxy_pass http://192.168.122.1:9091;
                    proxy_set_header Host $http_host;
                    proxy_set_header X-Forwarded-For $cip;
                    client_body_buffer_size 9M;
                    client_max_body_size 4000M;
                    proxy_max_temp_file_size 8m;
                    proxy_buffers 1024 4k;
                    proxy_read_timeout 300;
                    access_log /var/log/nginx/unode_acc.log main;
                }
            }`)

sitesMap.set("git.cloudhua.com", `
                location ~* /fop/convert {
                    proxy_buffers 2 4k;
                    proxy_busy_buffers_size 4k;
                    proxy_set_header X-Forwarded-Proto $scheme;
                    proxy_set_header X-Forwarded-Host  $http_host;
                    proxy_pass http://127.0.0.1:10002;
                    add_header Access-Control-Allow-Origin "*";
                }
            `)
sitesMap.set("docker.cloudhua.com", `
                location ~* /fop/imageprocessing {
                    proxy_busy_buffers_size 4k;
                    proxy_buffers 2 4k;
                    proxy_pass http://127.0.0.1:8800;
                    add_header Access-Control-Allow-Origin "*";
                }
            `)

function createSite(e) {
    var site = $('#new-site').val()
    if (!isValidDomain(site)) {
        alert("无效的域名\n示例:abc.xyz.com")
    }
    alert("OK:" + site)
}

function saveCfg() {
    var site = $('#currentsite').html()
    var siteData = $('#ngx-cfg-view').html()
    if (siteData.length === 0) {
        alert("站点内容为空")
        return
    }
    $.post('savesite?site=' + encodeURIComponent(site), siteData, function() {
        alert('success')
    }).fail(function(data) {
        alert(data['responseText'])
    })
}

function testCfg() {
    $.get('testsite', function(data) {
        alert(data)
    }).fail(function(data) {
        alert(data)
    })
}

function reloadCfg() {
    $.get('reloadsite', function(data) {
        console.log(data)
        alert(data)
    }).fail(function(data) {
        alert(data)
    })
}

function showCfg(site) {
    $.get('sitecontent?site=' + encodeURIComponent(site), function(data) {
        $('#ngx-cfg-view').html(data['content'])
    }).fail(function() {
        alert('fail')
    })
}

function isValidDomain(v) {
    if (typeof v !== 'string') return false

    var parts = v.split('.')
    if (parts.length <= 2) return false

    var tld = parts.pop()
    var tldRegex = /^[a-zA-Z0-9]+$/gi

    if (!tldRegex.test(tld)) return false

    var isValid = parts.every(function(host) {
        var hostRegex = /^(?!:\/\/)([a-zA-Z0-9]+|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])$/gi;
        return hostRegex.test(host)
    })

    return isValid
}


$(document).ready(function() {
    $('.ui.search').search({
        source: sites,
        onSelect: function(result, resp) {
            var site = result["title"]
            showCfg(site)
        }
    });

    $('#sites .item').on('click', function() {
            $(this).addClass('active').siblings().removeClass('active');
            var site = $(this).children(".site")[0].innerHTML
            $('#currentsite').html(site)
            showCfg(site)
        }
    );

    $('#create-site').on('click', createSite)
    $('#save-cfg').on('click', saveCfg)
    $('#test-cfg').on('click', testCfg)
    $('#reload-cfg').on('click', reloadCfg)
})
