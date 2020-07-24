getLogs();

$(document).ready(function() {

  //Determine id of Job/Log component and load respective jobs/logs
  var id = $(".STmain").attr("id");
  $(".STmain").find("#element").text(id);


  $(".STmain").find(".ltable").empty();

  //Loop over Logs JSON DB and add Logs related to parent summoning "id" (Bichito ID or "Hive")
  for (i = 0; i < logs.length; i++){
    var row = logs[i];
    if ((row.pid == id) || (row.chid == id)){
      $(".STmain").find(".ltable").append("<tr class=\"header\"><td>"+htmlencode.htmlEncode(row.time)+"</td><td colspan=\"4\">"+htmlencode.htmlEncode(row.error)+"</td></tr>");
    }
  }

})