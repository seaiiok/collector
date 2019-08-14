<!DOCTYPE html>

<head>
    <meta charset="UTF-8">
    <title>Collect</title>
</head>

<body>
    <style type="text/css">
        table.gridtable {
            font-family: verdana, arial, sans-serif;
            font-size: 11px;
            color: #303133;
            width: 80%
        }

        table.gridtable th {
            padding: 10px;
            color: #ffffff;
            background-color: #303133;
        }

        table.gridtable td {
            padding: 8px;
            background-color: #F2F6FC;
        }
    </style>

    <!-- Table goes JF room offline File Info by MES-->
    <table class="gridtable">
        <tr>
            <th>日 期</th>
            <th>文件地址</th>
            <th>文件路径</th>
            <th>状 态</th>
        </tr>

        {{ range $key,$value:=. }}
        <tr class="ctl">
            <td>{{$value.Date}}</td>
            <td>{{$value.IP}}</td>
            <td>{{$value.Path}}</td>
          
            <td id="{{$key}}" style="color: grey"></td>

            <script>
                var x = document.getElementById("{{$key}}")
                if ({{ $value.Status }} == "1") {
                    x.style.color = "green"
                    x.innerHTML = "完成"
                } else if ({{ $value.Status }} == "0") {
                    x.style.color = "red"
                    x.innerHTML = "未完成"
                }else {
                    x.innerHTML = "未知"
                }
            </script>

        </tr>
        {{ end }}
    </table>

</body>

</html>