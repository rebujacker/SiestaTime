
// Loads Stagings for the Attacks
function loadFormData(button) {

  if (stagings != null){
    for (i = 0; i < stagings.length; i++){
      var row = stagings[i];
      $("#stagingOpt").append("<option>"+row.name+"</option>");
    }
  }

}


$(document).ready(function() {
  //// Refresh on memory data and load it the sidetables for element creations
  getStagings();

  loadFormData();

  //Determine if elog is for hive,red or bichito, then load the logs for each case
  var name = $(".STmain").attr("id");
  $(".STmain").find(".element").text(name);

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
 
  var infoJson = JSON.parse(bichito.info)

  $("#pid").text(infoJson.pid.replace(/\n/g, ''));
  $("#arch").text(infoJson.arch.replace(/\n/g, ''));
  $("#os").text(infoJson.os.replace(/\n/g, ''));
  $("#osv").text(infoJson.osv.replace(/\n/g, ''));
  $("#hostname").text(infoJson.hostname.replace(/\n/g, ''));
  $("#mac").text(infoJson.mac.replace(/\n/g, ''));
  $("#buser").text(infoJson.user.replace(/\n/g, ''));
  $("#privileges").text(infoJson.privileges.replace(/\n/g, ''));
  $("#lastonline").text(bichito.lastchecked);


  $("#lastdomain").text(domain);

  
  $(".STmain").on("click","button",function() {
    var link = $(this);
    switch(link.attr("id")) {
      case "jobs":
        //console.log(link.attr("id"));
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


$("#submitInject").on('click',function(){

  var attack = $('#injectType').val();
  console.log(attack);
  //Create Job to send with two elements

  var createInjectJSON = {staging:$("#stagingOpt").val()}
  var data = {cid:"",jid:"",pid:"",chid:$(".STmain").attr("id"),job:attack,time:"",status:"",result:"",parameters:"["+JSON.stringify(createInjectJSON)+"]"};
  //data.push();
  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/job",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){
          console.log("Response Job:"+response[0].jid);
          if (response != null){
            console.log("Response Job:"+response[0].jid);
            return
          }
        }

    });

});