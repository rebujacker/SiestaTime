getLogs();

$(document).ready(function() {

  //Determine id of Job/Log component and load respective jobs/logs
  var id = $(".STmain").attr("id");
  $(".STmain").find("#element").text(id);

  //if (id == "Hive"){
    //for (i = 0; i < logs.length; i++){
      //var row = logs[i];
      //if (row.pid == "Hive"){
       // $(".STmain").find(".ltable").append("<tr class=\"header\"><td>"+row.time+"</td><td colspan=\"4\">"+row.error+"</td></tr>");
      //}
    //}
  //}else{
    $(".STmain").find(".ltable").empty();
    for (i = 0; i < logs.length; i++){
      var row = logs[i];
      if ((row.pid == id) || (row.chid == id)){
        $(".STmain").find(".ltable").append("<tr class=\"header\"><td>"+row.time+"</td><td colspan=\"4\">"+row.error+"</td></tr>");
      }
    }
  //}

})