/*

Forms.js Is the client code that will handle the actions to be performed in relation with Creation Forms

*/


// Loads VPS's and Domains for the creation of Implant list
// It Also loads the Existing implant list for creating Stagings
function loadFormDataDomains(button) {

  if (implants != null){
    for (i = 0; i < implants.length; i++){
      var row = implants[i];
      $("#implantOpt").append("<option>"+htmlencode.htmlEncode(row.name)+"</option>");
    }
  }

  if (vps != null){
    for (i = 0; i < vps.length; i++){
      var row = vps[i];
      $("#vpsOpt").append("<option>"+htmlencode.htmlEncode(row.name)+"</option>");
    }
  }

  if (domains != null){
    for (i = 0; i < domains.length; i++){
      var row = domains[i];
      if ((row.active == "No") && (row.dtype != 'gmail')){
        $("#domainOpt").append("<option>"+htmlencode.htmlEncode(row.name)+"</option>");
      }
    }
  }
}


function loadFormDataSaaS(button) {

  if (implants != null){
    for (i = 0; i < implants.length; i++){
      var row = implants[i];
      $("#implantOpt").append("<option>"+htmlencode.htmlEncode(row.name)+"</option>");
    }
  }

  if (vps != null){
    for (i = 0; i < vps.length; i++){
      var row = vps[i];
      $("#vpsOpt").append("<option>"+htmlencode.htmlEncode(row.name)+"</option>");
    }
  }

  if (domains != null){
    for (i = 0; i < domains.length; i++){
      var row = domains[i];
      if ((row.active == "No") && (row.dtype == 'gmail')){
        $("#domainOpt").append("<option>"+htmlencode.htmlEncode(row.name)+"</option>");
      }
    }
  }
}


//// Parameters options for one Network Module or other

//// This will be in charge of adding/remove redirectors set in the GUI for creating Implants
  

function getP(){
  p = $("#participants").val();
}


function addRow() {
  $(".participantRow").clone(true, true).appendTo("#participantTable");
  $("#redNumber").text($('#participantTable tr').length - 2);
}

function removeRow(button) {
  button.closest("tr").remove();
}

/* Doc ready */
$("#params").on('click','.add', function () {
    addRow();

  $(this).closest("tr").appendTo("#participantTable");
  if ($("#participantTable tr").length === 3) {
    $(".remove").hide();
  } else {
    $(".remove").show();
  }
});

$("#params").on('click','.remove',function () {

  if($("#participantTable tr").length === 3) {
    //alert("Can't remove row.");
    $(".remove").hide();
  } else if($("#participantTable tr").length - 1 ==3) {
    $(".remove").hide();
    removeRow($(this));

  } else {
    removeRow($(this));
  }
});


//Used to change Forms for different VPS types
$('#coms').change(function(){

  //// Refresh on memory data and load it the sidetables for element creations
  getVps();
  getDomains();
  getImplants();



  if ($('#coms').val() == 'gmailgo'){
    $("#params").empty();
    $("#params").append(`
    
    <div class="form-group">
        <label for="iname">Server/Red For SaaS Coms </label>
        <select id="vpsOpt" class="required-entry" name="vps"></select>      
    </div>
    
    <div>
      <table class="form-group" id="participantTable">
        <thead>
            <tr>
                <th>Redirector</th>
                <th>SaaS API Account</th>
            </tr>
        </thead>
        <tr class="participantRow" name="redirector">
            <td id="redNumber"></td>
            <td>
                <select id="domainOpt" class="required-entry" name="domain">
                </select>
            </td>
            <td><button class="btn btn-danger remove" type="button">Remove</button></td>
        </tr>
        <tr id="addButtonRow">
            <td colspan="4"><center><button class="btn btn-large btn-success add" type="button">Add</button></center></td>
        </tr>
      </table>
  </div>
  <button type="button" class="btn btn-primary" id="submitcreationImplantDomain">Create Implant</button>
    `);
    loadFormDataSaaS();
  }else{
    $("#params").empty();
    $("#params").append(`
    <div class="form-group">
      <label for="iname">TLS Port </label>
      <input type="text" class="form-control" name="comsparams" placeholder="">
    </div>

<div>
  <table class="form-group" id="participantTable">
        <thead>
            <tr>
                <th>Redirector</th>
                <th>VPS</th>
                <th>Domain</th>
            </tr>
        </thead>
        <tr class="participantRow" name="redirector">
            <td id="redNumber"></td>
            <td>
                <select id="vpsOpt" class="required-entry" name="vps">
                </select>
            </td>
            <td>
                <select id="domainOpt" class="required-entry" name="domain">
                </select>
            </td>
            <td><button class="btn btn-danger remove" type="button">Remove</button></td>
        </tr>
        <tr id="addButtonRow">
            <td colspan="4"><center><button class="btn btn-large btn-success add" type="button">Add</button></center></td>
        </tr>
  </table>
</div>
<button type="button" class="btn btn-primary" id="submitcreationImplantSaaS">Create Implant</button>
    `);
    loadFormDataDomains();
  }

  $("#redNumber").text($('#participantTable tr').length - 2);
  
});


// The Implant form Submission, will massage the data to adapt itself to the Hive createImplant Format:

/*
type CreateImplant struct {
    Name string   `json:"name"`
    Ttl string   `json:"ttl"`
    Resptime string   `json:"resptime"`
    Coms string   `json:"coms"`
    ComsParams string `json:"comsparams"`
    Persistence string `json:"persistence"`
    Redirectors  []Red `json:"redirectors"`
}
*/
$("#params").on('click','#submitcreationImplantDomain',function () {

  //Serialize form in the correct way

  //Transform the array in one JSON STRING
  function objectifyForm(formArray) {
    var returnArray = {};
    var arrayRedirectors = [];
    var vps = "";
    for (var i = 0; i < formArray.length; i++){
      if (formArray[i]['name'] == 'vps'){
        vps = formArray[i]['value']
        returnArray[formArray[i]['name']] = formArray[i]['value'];
      }else if (formArray[i]['name'] == 'domain'){
        var tempObject = {};
        tempObject['vps'] = vps;
        tempObject['domain'] = formArray[i]['value'];
        arrayRedirectors.push(tempObject)
      }else{
        returnArray[formArray[i]['name']] = formArray[i]['value'];
      }
    }

    returnArray['redirectors'] = arrayRedirectors
    return returnArray;
  }

  var createImplantJSON = objectifyForm($("#createimplantform").serializeArray());
  console.log(objectifyForm($("#createimplantform").serializeArray()));


  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"createImplant",time:"",status:"",result:"",parameters:"["+JSON.stringify(createImplantJSON)+"]"};
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



$("#params").on('click','#submitcreationImplantSaaS',function () {

  //Serialize form in the correct way

  //Transform the array in one JSON STRING
  function objectifyForm(formArray) {
    var returnArray = {};
    var arrayRedirectors = [];
    for (var i = 0; i < formArray.length; i++){
      if (formArray[i]['name'] == 'vps'){
        var tempObject = {};
        tempObject[formArray[i]['name']] = formArray[i]['value'];
        i++;
        tempObject[formArray[i]['name']] = formArray[i]['value'];
        arrayRedirectors.push(tempObject)
      }else{
        returnArray[formArray[i]['name']] = formArray[i]['value'];
      }
    }

    returnArray['redirectors'] = arrayRedirectors
    return returnArray;
  }

  var createImplantJSON = objectifyForm($("#createimplantform").serializeArray());
  console.log(objectifyForm($("#createimplantform").serializeArray()));


  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"createImplant",time:"",status:"",result:"",parameters:"["+JSON.stringify(createImplantJSON)+"]"};
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

//// Vps Creation Form: Parameters and function

//Used to change Forms for different VPS types
$('#vtype').change(function(){

  switch($('#vtype').val()) {
    case 'aws_instance':
      $("#params").empty();
      $("#params").append(`
    <form role="form" id="vpsparamsform">
    <div class="form-group">
      <label for="iname"> Access Key </label>
      <input type="text" class="form-control" name="accesskey" id="iname" placeholder="VPS Access Key...">
    </div>

    <div class="form-group">
      <label for="iname"> Secret Key </label>
      <input type="text" class="form-control" name="secretkey" id="iname" placeholder="VPS Secret Key...">
    </div>

    <div class="form-group">
      <label for="iname"> Region </label>
      <input type="text" class="form-control" name="region" id="iname" placeholder="VPS Region...">
    </div>

    <div class="form-group">
      <label for="iname"> AMI </label>
      <input type="text" class="form-control" name="ami" id="iname" placeholder="VPS AMI...">
    </div>

    <div class="form-group">
      <label for="iname"> SSH Keyname </label>
      <input type="text" class="form-control" name="sshkeyname" id="iname" placeholder="VPS SSH Keyname...">
    </div>

    <div class="form-group">
      <label for="iname"> SSH Key </label>
      <textarea class="resizable_textarea" name="sshkey" rows="10" cols="30" placeholder="VPS SSH PEM Key..."></textarea> 
    </div>
    </form>   
    `);
      break;
    case 'azure_instance':
      $("#params").empty();
      $("#params").append(`
        <h1>Azure!</h1>  
      `); 
      break; 
  }
});

$("#submitcreationvps").on('click',function(){

  //TO-DO: Serializing logic for variable parameter field...

  //Transform the array in one JSON STRING
  // [{}]
  function objectifyForm(formArrayVps,formArrayParameters) {
    var returnArray = {};
    var arrayParameters = {};
    for (var i = 0; i < formArrayVps.length; i++){
      if (formArrayVps[i]['name'] == 'name'){
        returnArray[formArrayVps[i]['name']] = formArrayVps[i]['value'];
      }
      if (formArrayVps[i]['name'] == 'vtype'){
        returnArray[formArrayVps[i]['name']] = formArrayVps[i]['value'];
        for (var y = 0; y < formArrayParameters.length; y++){
          arrayParameters[formArrayParameters[y]['name']] = formArrayParameters[y]['value'];
        }
      }
    }
    returnArray['parameters'] = JSON.stringify(arrayParameters);
    return returnArray;
  }

  //Serialize form in the correct way
  var createVpsJSON = objectifyForm($("#createvpsform").serializeArray(),$("#vpsparamsform").serializeArray());

  console.log(JSON.stringify(createVpsJSON));

  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"createVPS",time:"",status:"",result:"",parameters:"["+JSON.stringify(createVpsJSON)+"]"};
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



//// Domain Creation Form: Parameters and function

//Used to change Forms for different VPS types
$('#dtype').change(function(){

  switch($('#dtype').val()) {
    case 'godaddy':
      $("#params").empty();
      $("#params").append(`
    <form role="form" id="domainparamsform">
    
    <div class="form-group">
      <label for="iname"> Domain </label>
      <input type="text" class="form-control" name="domain" id="iname" placeholder="domain.xzy...">
    </div>
    
    <div class="form-group">
      <label for="iname"> Access Key </label>
      <input type="text" class="form-control" name="domainkey" id="iname" placeholder="Domain Access Key...">
    </div>

    <div class="form-group">
      <label for="iname"> Secret Key </label>
      <input type="text" class="form-control" name="domainsecret" id="iname" placeholder="Domain Secret Key...">
    </div>
    </form>
    <button type="button" class="btn btn-primary" id="submitcreationdomain">Create Domain</button>   
    `);
      break;
    case 'gmail':
      $("#params").empty();
      $("#params").append(`
    <form role="form" id="domainparamsform">
    <div class="form-group">
      <label for="iname"> Credentials.json </label>
      <textarea class="resizable_textarea" name="creds" cols="30" placeholder="Gmail Cred Json File..."></textarea>
    </div>

    <div class="form-group">
      <label for="iname"> Token.json </label>
      <textarea class="resizable_textarea" name="token" cols="30" placeholder="Gmail Access/Refresh Token..."></textarea>
    </div>
    </form>
    <button type="button" class="btn btn-primary" id="submitcreationdomainSaaS">Create Domain</button>    
      `); 
      break; 
  }
});

$("#params").on('click','#submitcreationdomain',function(){

  //TO-DO: Serializing logic for variable parameter field...

  //Transform the array in one JSON STRING
  // [{}]
  function objectifyForm(formArrayDomain,formArrayParameters) {
    var returnArray = {};
    var arrayParameters = {};
    for (var i = 0; i < formArrayDomain.length; i++){
      if (formArrayDomain[i]['name'] == 'name'){
        returnArray[formArrayDomain[i]['name']] = formArrayDomain[i]['value'];
      }

      if (formArrayDomain[i]['name'] == 'dtype'){
        returnArray[formArrayDomain[i]['name']] = formArrayDomain[i]['value'];
        for (var y = 0; y < formArrayParameters.length; y++){
          if (formArrayParameters[y]['name'] == 'domain'){
            returnArray['domain'] =  formArrayParameters[y]['value'];
          }
          arrayParameters[formArrayParameters[y]['name']] = formArrayParameters[y]['value'];
        }
      }
    }

    returnArray['parameters'] = JSON.stringify(arrayParameters);
    return returnArray;
  }

  //Serialize form in the correct way

  var createDomainJSON = objectifyForm($("#createdomainform").serializeArray(),$("#domainparamsform").serializeArray());

  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"createDomain",time:"",status:"",result:"",parameters:"["+JSON.stringify(createDomainJSON)+"]"};
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

$("#params").on('click','#submitcreationdomainSaaS',function(){

  //TO-DO: Serializing logic for variable parameter field...

  //Transform the array in one JSON STRING
  // [{}]
  function objectifyForm(formArrayDomain,formArrayParameters) {
    var returnArray = {};
    var arrayParameters = {};
    for (var i = 0; i < formArrayDomain.length; i++){
      if (formArrayDomain[i]['name'] == 'name'){
        returnArray[formArrayDomain[i]['name']] = formArrayDomain[i]['value'];
        returnArray['domain'] =  formArrayDomain[i]['value'];
      }

      if (formArrayDomain[i]['name'] == 'dtype'){
        returnArray[formArrayDomain[i]['name']] = formArrayDomain[i]['value'];
        for (var y = 0; y < formArrayParameters.length; y++){
          arrayParameters[formArrayParameters[y]['name']] = formArrayParameters[y]['value'];
        }
      }
    }

    returnArray['parameters'] = JSON.stringify(arrayParameters);
    return returnArray;
  }

  //Serialize form in the correct way

  var createDomainJSON = objectifyForm($("#createdomainform").serializeArray(),$("#domainparamsform").serializeArray());

  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"createDomain",time:"",status:"",result:"",parameters:"["+JSON.stringify(createDomainJSON)+"]"};
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



//Used to change Forms for different staging types
$('#stype').change(function(){

  switch($('#stype').val()) {
    case 'https_droplet_letsencrypt':
      $("#params").empty();
      $("#params").append(`
    <form role="form" id="stagingparamsform">
    <div class="form-group">
      <label for="iname">Droplet HTTPs Port </label>
      <input type="text" class="form-control" name="httpsport" placeholder="1244">
    </div>

    <div class="form-group">
      <label for="iname">Path for Implants</label>
      <input type="text" class="form-control" name="path" placeholder="namepath">
    </div>
    </form>   
    `);
      break;
      loadFormDataDomains();
    case 'https_msft_letsencrypt':
      $("#params").empty();
      $("#params").append(`
    <form role="form" id="stagingparamsform">
    <div class="form-group">
      <label for="iname">MSFT HTTPs Port </label>
      <input type="text" class="form-control" name="httpsport" placeholder="1244">
    </div>
    </form> 
      `); 
      break; 
      loadFormDataDomains();
    case 'https_empire_letsencrypt':
      $("#params").empty();
      $("#params").append(`
    <form role="form" id="stagingparamsform">
    <div class="form-group">
      <label for="iname">Empire HTTPs Port </label>
      <input type="text" class="form-control" name="httpsport" placeholder="1244">
    </div>
    </form> 
      `); 
      loadFormDataDomains();
      break; 
  }
});



$("#submitcreationstaging").on('click',function(){

  //TO-DO: Serializing logic for variable parameter field...

  //Transform the array in one JSON STRING
  // [{}]
  function objectifyForm(formArrayStaging,formArrayParameters) {
    var returnArray = {};
    var arrayParameters = {};
    for (var i = 0; i < formArrayStaging.length; i++){
      if (formArrayStaging[i]['name'] == 'stype'){
        returnArray[formArrayStaging[i]['name']] = formArrayStaging[i]['value'];
        for (var y = 0; y < formArrayParameters.length; y++){
          arrayParameters[formArrayParameters[y]['name']] = formArrayParameters[y]['value'];
        }
      }else{
        returnArray[formArrayStaging[i]['name']] = formArrayStaging[i]['value'];
      }
    }
    returnArray['parameters'] = JSON.stringify(arrayParameters);
    return returnArray;
  }

  //Serialize form in the correct way

  var createStagingJSON = objectifyForm($("#createstagingform").serializeArray(),$("#stagingparamsform").serializeArray());
  console.log($("#createstagingform").serializeArray());
  
  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"createStaging",time:"",status:"",result:"",parameters:"["+JSON.stringify(createStagingJSON)+"]"};
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



$("#submitcreationreport").on('click',function(){

  //TO-DO: Serializing logic for variable parameter field...

  //Transform the array in one JSON STRING
  // [{}]
  function objectifyForm(formArrayVps) {
    var returnArray = {};
    for (var i = 0; i < formArrayVps.length; i++){

        returnArray[formArrayVps[i]['name']] = formArrayVps[i]['value'];

    }

    return returnArray;
  }

  //Serialize form in the correct way

  var createVpsJSON = objectifyForm($("#createreportform").serializeArray());


  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"createReport",time:"",status:"",result:"",parameters:"["+JSON.stringify(createVpsJSON)+"]"};
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




