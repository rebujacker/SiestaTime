
// Loads Stagings for the Attacks
function loadFormData(button) {

  if (stagings != null){
    for (i = 0; i < stagings.length; i++){
      var row = stagings[i];
      $("#stagingOpt").append("<option>"+htmlencode.htmlEncode(row.name)+"</option>");
    }
  }

}


$(document).ready(function() {
  //// Refresh on memory data and load it the sidetables for element creations
  getStagings();
  loadFormData();

  //Determine if elog is for hive,red or bichito, then load the logs for each case
  var name = $(".STmain").attr("id");
  $(".STmain").find(".element").text(htmlencode.htmlEncode(name));

  for (i = 0; i < bichitos.length; i++){
    if (bichitos[i].bid == name){
      var bichito = bichitos[i];
    }
  }

  for (i = 0; i < redirectors.length; i++){
    if (redirectors[i].rid == bichito.rid){
      var domain = redirectors[i].domainame;
    }
  }
 
 //Fetch Bichito Information to panel
  var infoJson = JSON.parse(bichito.info)

  $("#pid").text(htmlencode.htmlEncode(infoJson.pid.replace(/\n/g, '')));
  $("#arch").text(htmlencode.htmlEncode(infoJson.arch.replace(/\n/g, '')));
  $("#os").text(htmlencode.htmlEncode(infoJson.os.replace(/\n/g, '')));
  $("#osv").text(htmlencode.htmlEncode(infoJson.osv.replace(/\n/g, '')));
  $("#hostname").text(htmlencode.htmlEncode(infoJson.hostname.replace(/\n/g, '')));
  $("#mac").text(htmlencode.htmlEncode(infoJson.mac.replace(/\n/g, '')));
  $("#buser").text(htmlencode.htmlEncode(infoJson.user.replace(/\n/g, '')));
  $("#privileges").text(htmlencode.htmlEncode(infoJson.privileges.replace(/\n/g, '')));
  $("#lastonline").text(htmlencode.htmlEncode(bichito.lastchecked));
  $("#status").text(htmlencode.htmlEncode(bichito.status));


  $("#lastdomain").text(domain);

  //Load components options: Jobs,logs, console
  $(".btn").unbind().click(function() {
    console.log("happening")
    var link = $(this);
    switch(link.attr("id")) {
      case "jobs":
        $(".STmain").load('./components/jobs/jobs.html')
        break;
      case "logs":
        $(".STmain").load('./components/logs/logs.html')
        break;
      case "console":
        $(".STmain").find("#binteraction").load('./components/console/console.html')
        break;
      default:
    }
  });

});


//Similarly to forms.js, inject a set of params in function of the type of injection
$('#injectType').change(function(){

  //// Refresh on memory data and load it the sidetables for element creations
  getStagings();



  if ($('#injectType').val() == 'injectRevSshShellOffline'){
    $("#injectParams").empty();
    $("#injectParams").append(`
    
      <div class="form-group">
        <label for="iname">Domain or IP - SSH Server </label>
        <input type="text" class="form-control" name="domain" placeholder="Domain/IP...">
      </div>

      <div class="form-group">
        <label for="timeouts">Port </label>
        <input type="text" class="form-control" name="port" placeholder="ssh port...">
      </div>
      <div class="form-group">
        <label for="timeouts">UserName of target SSH </label>
        <input type="text" class="form-control" name="user" placeholder="ssh server username...">
      </div>
      <div class="form-group">
        <label for="timeouts">SSH KEY </label>
        <textarea class="resizable_textarea" name="sshkey" rows="10" cols="30" placeholder="SSH PEM Key..."></textarea> 
      </div>
      <button type="button" class="btn btn-primary" id="submitInjectOffline">Inject</button>
    `);

  }else {
    $("#injectParams").empty();
    $("#injectParams").append(`
    
      <div>
        <table class="table table-striped table-bordered">
          <thead>
            <tr>
              <th>Staging Server</th>
            </tr>
           </thead>
           <tr name="stagings">
              <td>
              <select id="stagingOpt" class="form-control" required>
              </select>
              </td>
            </tr>
       </table>
      </div>
      <button type="button" class="btn btn-primary" id="submitInject">Inject</button>
    `);  
    loadFormData(); 
  }

});




/*This Job will respect the following JSON Structure on "parameters":
type InjectEmpire struct {
    Staging string   `json:"staging"`
}
*/
$("#injectParams").on('click','#submitInject',function(){

  var attack = $('#injectType').val();

  var createInjectJSON = {staging:$("#stagingOpt").val()}
  var data = {cid:"",jid:"",pid:"",chid:$(".STmain").attr("id"),job:attack,time:"",status:"",result:"",parameters:"["+JSON.stringify(createInjectJSON)+"]"};

  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/job",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){
          if (response != null){
            return
          }
        }

    });

});

//Special Command for Offline
$("#injectParams").on('click','#submitInjectOffline',function(){


  function objectifySimpleForm(formArray) {
    var returnArray = {};

    for (var i = 0; i < formArray.length; i++){
        returnArray[formArray[i]['name']] = formArray[i]['value']; 
    }

    return returnArray;
  }

  var attack = $('#injectType').val();
  var createInjectJSON = objectifySimpleForm($("#createInjectform").serializeArray());

  var data = {cid:"",jid:"",pid:"",chid:$(".STmain").attr("id"),job:attack,time:"",status:"",result:"",parameters:"["+JSON.stringify(createInjectJSON)+"]"};

  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/job",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){

          if (response != null){
            return
          }
        }

    });

});