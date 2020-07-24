
$(document).ready(function() {


  //Determine if elog is for hive,red or bichito, then load the logs for each case
  var name = $(".STmain").attr("id");

  $(".element").text(htmlencode.htmlEncode(name))
  $(".element").attr("value",htmlencode.htmlEncode(name));

})

/* Craft a Job with the following JSON Object towards client:

type ReportObject struct {
    Name string   `json:"name"`
}
*/

$("#submitdownloadreport").on('click',function() {

  var reportName = $(".element").attr("value");

  //Create Job to send with two elements
  var data = {name:reportName};

  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/report",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){

        }

    });

});