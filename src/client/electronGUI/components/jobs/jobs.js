getJobs();

$(document).ready(function() {


  //Determine id of Job/Log component and load respective jobs/logs
  var id = $(".STmain").attr("id");
  $(".STmain").find(".element").text(id);

  //if (id == "Hive"){
    //for (i = 0; i < jobs.length; i++){
      //var row = jobs[i];
      //if (row.pid == "Hive"){
        //$(".STmain").find(".jtable").append("<tr class=\"header\"><td>"+row.cid+"</td><td>"+row.jid+"</td><td>"+row.time+"</td><td>"+row.job+"</td><td>"+row.status+"</td></tr><tr style=\"display: none;\"><td colspan=\"5\">"+row.parameters+"</td><tr style=\"display: none;\"><td colspan=\"5\">"+row.result+"</td></tr></tr>");
      //}
    //}
  //}else{
    $(".STmain").find(".jtable").empty();
    for (i = 0; i < jobs.length; i++){
      var row = jobs[i];
      if ((row.pid == id) || (row.chid == id)){
        if (row.result.length >= 10000) {
          row.result = "Too Large Output - blob";
        }
        $(".STmain").find(".jtable").append("<tr class=\"header\"><td>"+htmlencode.htmlEncode(row.cid)+"</td><td>"+htmlencode.htmlEncode(row.jid)+"</td><td>"+htmlencode.htmlEncode(row.time)+"</td><td>"+htmlencode.htmlEncode(row.job)+"</td><td>"+htmlencode.htmlEncode(row.status)+"</td></tr><tr style=\"display: none;\"><td colspan=\"5\">"+htmlencode.htmlEncode(row.parameters)+"</td><tr style=\"display: none;\"><td colspan=\"5\">"+htmlencode.htmlEncode(row.result)+"</td></tr></tr>");
      }
    }
  //}


  if ($(".STmain").find('.table').length > 0) {
    $(".STmain").find('.table .header').on("click", function() {
      
      $(this).toggleClass("active", "").nextUntil('.header').css('display', function(i, v) {
        return this.style.display === 'table-row' ? 'none' : 'table-row';
      });
    });
  }
  

})