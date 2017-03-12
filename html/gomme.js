var uri = '/gomme-api/'; // 'http://192.168.2.2:8080/'
var cmd = new Array();

function send_cmd(cmd_name) {
  if ( cmd[cmd_name] == "" ) {
    console.log("Command not found ", cmd_name);
    return;
  }
  $.ajax({
    url: cmd[cmd_name],
  })
  .done(function( data ) {
    if ( console && console.log ) {
      console.log(data);
    }
  });
}
{{range .Buttons}}
$('#{{.}}').on('click', function(v) { send_cmd ('{{.}}'); } )
cmd['{{.}}'] = uri + '{{.}}';
{{end}}
/*
$('#ex1').slider({
	formatter: function(value) {
		return 'Current value: ' + value;
	}
});
*/
