/*

Forms.js Is the client code that will handle the actions to be performed in relation with Creation Forms

*/

//Refresh key data on creation
  getVps();
  getDomains();
  getImplants();
  getStagings();

// Loads VPS's and Domains for the creation of Implant list
// It Also loads the Existing implant list for creating Stagings
function loadFormDataDomains() {

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


function loadFormDataSaaS() {

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

loadFormDataDomains();
loadFormDataSaaS();

//// Parameters options for one Network Module or other

//// This will be in charge of adding/remove redirectors set in the GUI for creating Implants
  

function getP(){
  p = $("#participants").val();
}


function addRow() {
  $(".participantRow:first").clone(true, true).appendTo("#participantTable:last");
  $("#redNumber").text("Redirectors: " + ($('#participantTable tr').length - 2));
}

function removeRow(button) {
  button.closest("tr").remove();
  $("#redNumber").text("Redirectors: " + ($('#participantTable tr').length - 2));
}

/* Doc ready */
$("#netParams").on('click','.add', function () {
    addRow();

  $(this).closest("tr").appendTo("#participantTable");
  if ($("#participantTable tr").length === 3) {
    $(".remove").hide();
  } else {
    $(".remove").show();
  }
});

$("#netParams").on('click','.remove',function () {

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


//Userland Persistence Module Options
$('#persistenceosx').change(function(){

  if ($('#persistenceosx').val() == 'launchd'){
    $("#userlandPersistenceOSXParams").empty();
    $("#userlandPersistenceOSXParams").append(`
    
    <form role="form" id="userlandpersistenceosxparamsform">
      <div class="form-group">
        <label for="launchdname"> Launchd Name (~/Library/LaunchAgents/com.name.agent.plist)</label>
        <input type="text" class="form-control" name="launchdname" id="launchdname" placeholder="Name...">
      </div>

      <div class="form-group">
        <label for="implantpath"> Implant Path (Relative to Default User Home Folder) </label>
        <input type="text" class="form-control" name="implantpath" id="implantpath" placeholder="folder/folder/filename...">
      </div>

    </form>  
    `);
  }else{
    $("#userlandPersistenceOSXParams").empty();
  }

});

$('#persistencewindows').change(function(){

  if ($('#persistencewindows').val() == 'schtasks'){
    $("#userlandPersistenceWindowsParams").empty();
    $("#userlandPersistenceWindowsParams").append(`
    
    <form role="form" id="userlandpersistencewindowsparamsform">
      <div class="form-group">
        <label for="taskname"> Schtask Name </label>
        <input type="text" class="form-control" name="taskname" id="taskname" placeholder="Name...">
      </div>

      <div class="form-group">
        <label for="implantpath"> Implant Path (Relative to Default User Home Folder) </label>
        <input type="text" class="form-control" name="implantpath" id="implantpath" placeholder="folder\\folder\\filename...">
      </div>

    </form>  
    `);
  }else{
    $("#userlandPersistenceWindowsParams").empty();
  }

});


$('#persistencelinux').change(function(){

  if ($('#persistencelinux').val() == 'linuxautostart'){
    $("#userlandPersistenceLinuxParams").empty();
    $("#userlandPersistenceLinuxParams").append(`
    
    <form role="form" id="userlandpersistencelinuxparamsform">
      <div class="form-group">
        <label for="cronname"> Autostart File Name ($XDG_CONFIG_HOME/.config/autostart/name.desktop) </label>
        <input type="text" class="form-control" name="autostartname" id="autostartname" placeholder="Name...">
      </div>

      <div class="form-group">
        <label for="implantpath"> Implant Path (Relative to Default User Home Folder) </label>
        <input type="text" class="form-control" name="implantpath" id="implantpath" placeholder="folder/folder/filename...">
      </div>

    </form>  
    `);
  }else{
    $("#userlandPersistenceLinuxParams").empty();
  }

});



//Used to change Forms for different VPS types
$('#coms').change(function(){

  //// Refresh on memory data and load it the sidetables for element creations
  getVps();
  getDomains();
  getImplants();



  if ($('#coms').val() == 'selfsignedhttpsgo'){
    $("#netParams").empty();
    $("#netParams").append(`
    <div class="form-group">
      <label for="iname">TLS Port </label>
      <input type="text" class="form-control" name="comsparams" placeholder="">
    </div>

<div>
  <label id="redNumber"> Redirectors: 1 </label>
  <table class="form-group" id="participantTable">
        <thead>
            <tr>
                <th>VPS</th>
                <th>Domain</th>
            </tr>
        </thead>
        <tr class="participantRow" name="redirector">
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
<button type="button" class="btn btn-primary" id="submitcreationImplantDomain">Create Implant</button>
    `);
    loadFormDataDomains();
  }else if ($('#coms').val() == 'paranoidhttpsgo'){
    $("#netParams").empty();
    $("#netParams").append(`
    <div class="form-group">
      <label for="iname">TLS Port </label>
      <input type="text" class="form-control" name="comsparams" placeholder="">
    </div>

<div>
  <label id="redNumber"> Redirectors: 1 </label>
  <table class="form-group" id="participantTable">
        <thead>
            <tr>
                <th>VPS</th>
                <th>Domain</th>
            </tr>
        </thead>
        <tr class="participantRow" name="redirector">
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
<button type="button" class="btn btn-primary" id="submitcreationImplantDomain">Create Implant</button>
    `);
    loadFormDataDomains();

  }else if ($('#coms').val() == 'gmailgo'){
    $("#netParams").empty();
    $("#netParams").append(`
    
    <div class="form-group">
        <label for="iname">Server/Red For SaaS Coms </label>
        <select id="vpsOpt" class="required-entry" name="vps"></select>      
    </div>
    
    <div>
      <label id="redNumber"> Redirectors: 1 </label>
      <table class="form-group" id="participantTable">
        <thead>
            <tr>
                <th>SaaS API Account</th>
            </tr>
        </thead>
        <tr class="participantRow" name="redirector">
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
    loadFormDataSaaS();
  }else if ($('#coms').val() == 'gmailmimic'){
    $("#netParams").empty();
    $("#netParams").append(`
 
    <div class="form-group">
        <label for="iname">User Agent </label>
        <input type="text" class="form-control" id="comsparam1" name="comsparam1" placeholder="Mozilla/5.0 (X11; Linux x86_64) AppleWeb...">     
    </div>

    <div class="form-group">
        <label for="iname">TLS Fingenprint (JA3 provided)</label>
        <input type="text" class="form-control" name="comsparams" id="comsparam2" name="comsparam2" placeholder="71,4865-4866-4867-49195-491...">     
    </div>

    <div class="form-group">
        <label for="iname">Server/Red For SaaS Coms </label>
        <select id="vpsOpt" class="required-entry" name="vps"></select>      
    </div>
    
    <div>
      <label id="redNumber"> Redirectors: 1 </label>
      <table class="form-group" id="participantTable">
        <thead>
            <tr>
                <th>SaaS API Account</th>
            </tr>
        </thead>
        <tr class="participantRow" name="redirector">
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
    loadFormDataSaaS();


  
  }else if ($('#coms').val() == 'selfsignedhttpsgoOffline'){
    $("#netParams").empty();
    $("#netParams").append(`
    <div class="form-group">
      <label for="iname">TLS Port </label>
      <input type="text" class="form-control" name="comsparams" placeholder="">
    </div>

<div>
  <label id="redNumber"> Redirectors: 1 </label>
  <table class="form-group" id="participantTable">
        <thead>
            <tr>
                <th>Domain or IP</th>
            </tr>
        </thead>
        <tr class="participantRow" name="redirector">
            <td>
                <input type="hidden" name="vps" />
            </td>
            <td>
                <input type="text" name="domain" />
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
    loadFormDataDomains();
  
  }else if ($('#coms').val() == 'paranoidhttpsgoOffline'){
    $("#netParams").empty();
    $("#netParams").append(`
    <div class="form-group">
      <label for="iname">TLS Port </label>
      <input type="text" class="form-control" name="comsparams" placeholder="">
    </div>

<div>
  <label id="redNumber"> Redirectors: 1 </label>
  <table class="form-group" id="participantTable">
        <thead>
            <tr>
                <th>Domain</th>
            </tr>
        </thead>
        <tr class="participantRow" name="redirector">
            <td>
                <input type="hidden" name="vps" />
            </td>
            <td>
                <input type="text" name="domain" />
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
    loadFormDataDomains();
  
  }else if ($('#coms').val() == 'gmailgoOffline'){
    $("#netParams").empty();
    $("#netParams").append(`
    
    <div class="form-group">
        <input type="hidden" name="vps" />     
    </div>
    
    <div>
      <label id="redNumber"> Redirectors: 1 </label>
      <table class="form-group" id="participantTable">
        <thead>
            <tr>
                <th>SaaS API Account</th>
            </tr>
        </thead>
        <tr class="participantRow" name="redirector">
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
    loadFormDataSaaS();
  
  }else if ($('#coms').val() == 'gmailmimicOffline'){
    $("#netParams").empty();
    $("#netParams").append(`
 
    <div class="form-group">
        <label for="iname">User Agent </label>
        <input type="text" class="form-control" id="comsparam1" name="comsparam1" placeholder="Mozilla/5.0 (X11; Linux x86_64) AppleWeb...">     
    </div>

    <div class="form-group">
        <label for="iname">TLS Fingenprint (JA3 provided)</label>
        <input type="text" class="form-control" name="comsparams" id="comsparam2" name="comsparam2" placeholder="71,4865-4866-4867-49195-491...">     
    </div>

    <div class="form-group">
        <input type="hidden" name="vps" />     
    </div>
    
    <div>
      <label id="redNumber"> Redirectors: 1 </label>
      <table class="form-group" id="participantTable">
        <thead>
            <tr>
                <th>SaaS API Account</th>
            </tr>
        </thead>
        <tr class="participantRow" name="redirector">
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
    loadFormDataSaaS();
  }
  
});


// The Implant form Submission, will massage the data to adapt itself to the Hive createImplant Format:

/*This Job will respect the following JSON Structure on "parameters":
type CreateImplant struct {
    Offline string   `json:"offline"`
    Name string   `json:"name"`
    Ttl string   `json:"ttl"`
    Resptime string   `json:"resptime"`
    Coms string   `json:"coms"`
    ComsParams string `json:"comsparams"`
    PersistenceOsx string `json:"persistenceosx"`
    PersistenceOsxP string `json:"persistenceosxp"`
    PersistenceWindows string `json:"persistencewindows"`
    PersistenceWindowsP string `json:"persistencewindowsp"`
    PersistenceLin string `json:"persistencelin"`
    PersistenceLinP string `json:"persistencelinp"`
    Redirectors  []Red `json:"redirectors"` //Array of Serialized "Red" JSON Objects Strings
}

type Red struct{
    Vps string `json:"vps"`
    Domain string `json:"domain"`
}

*/
$("#netParams").on('click','#submitcreationImplantDomain',function () {

  
  //Simple HTML-JSON Serialization: input1,input1 to {key:"value"}
  function objectifySimpleForm(formArray) {
    var returnArray = {};

    for (var i = 0; i < formArray.length; i++){
        returnArray[formArray[i]['name']] = formArray[i]['value']; 
    }

    return returnArray;
  }


  //Custom HTML-JSON Serialization: input1,input1 to {key:"value"}
  function objectifyImplantForm(formArray) {
    var returnArray = {};
    var arrayComsParam = [];
    var arrayRedirectors = [];

    //Loop over the provided array
    var vps = "";
    for (var i = 0; i < formArray.length; i++){

      //Create a new inner array for the selected redirectors
      if (formArray[i]['name'] == 'vps'){
        vps = formArray[i]['value']
      }else if (formArray[i]['name'] == 'domain'){
        var tempObject = {};
        tempObject['vps'] = vps;
        tempObject['domain'] = formArray[i]['value'];
        arrayRedirectors.push(tempObject)

      //ComsParams Array
      }else if (formArray[i]['name'].startsWith('comsparam')){
        var tempObject = {};
        arrayComsParam.push(formArray[i]['value'])

      }else{
        returnArray[formArray[i]['name']] = formArray[i]['value'];
      }
    }

    //In relation to the coms params, determine if it is an Offline Implant or not
    if (returnArray['coms'] == "selfsignedhttpsgoOffline") {
      returnArray['coms'] = "selfsignedhttpsgo"
      returnArray['offline'] = "Yes"
    }else if (returnArray['coms'] == "paranoidhttpsgoOffline"){

      returnArray['coms'] = "paranoidhttpsgo"
      returnArray['offline'] = "Yes"
    }else{
      returnArray['offline'] = "No"
    }
    
    
    returnArray['comsparams'] = arrayComsParam
    returnArray['redirectors'] = arrayRedirectors
    returnArray['persistenceosxp'] = JSON.stringify(objectifySimpleForm($("#userlandpersistenceosxparamsform").serializeArray()));
    returnArray['persistencewindowsp'] = JSON.stringify(objectifySimpleForm($("#userlandpersistencewindowsparamsform").serializeArray()));
    returnArray['persistencelinuxp'] = JSON.stringify(objectifySimpleForm($("#userlandpersistencelinuxparamsform").serializeArray()));

    return returnArray;
  }

  var createImplantJSON = objectifyImplantForm($("#createimplantform").serializeArray());
  ////console.log(createImplantJSON);


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
          if (response != null){
            return
          }
        }

    });

});


// When the network module selected is SaaS , function to send commands need to change because redirections have a different massaging
$("#netParams").on('click','#submitcreationImplantSaaS',function () {

  //Serialize form in the correct way

  //Transform the array in one JSON STRING
  function objectifyForm(formArray) {
    var returnArray = {};
    var arrayComsParam = [];
    var arrayRedirectors = [];

    for (var i = 0; i < formArray.length; i++){
      if (formArray[i]['name'] == 'vps'){
        var tempObject = {};
        tempObject[formArray[i]['name']] = formArray[i]['value'];
        i++;
        tempObject[formArray[i]['name']] = formArray[i]['value'];
        arrayRedirectors.push(tempObject)

      //ComsParams Array
      }else if (formArray[i]['name'].startsWith('comsparam')){
        var tempObject = {};
        arrayComsParam.push(formArray[i]['value'])
        
      }else{
        returnArray[formArray[i]['name']] = formArray[i]['value'];
      }
    }

    if (returnArray['coms'] == "gmailgoOffline") {

      returnArray['coms'] = "gmailgo"
      returnArray['offline'] = "Yes"
    }else if (returnArray['coms'] == "gmailmimicOffline"){

      returnArray['coms'] = "gmailmimic"
      returnArray['offline'] = "Yes"
    }else{
      returnArray['offline'] = "No"
    }

    returnArray['comsparams'] = arrayComsParam
    returnArray['redirectors'] = arrayRedirectors
    return returnArray;
  }

  var createImplantJSON = objectifyForm($("#createimplantform").serializeArray());

  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"createImplant",time:"",status:"",result:"",parameters:"["+JSON.stringify(createImplantJSON)+"]"};
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


/*This Job will respect the following JSON Structure on "parameters":
type Vps struct {
    Name string   `json:"name"`
    Vtype  string   `json:"vtype"`          // Aamazon, Azure, Lineage...
    Parameters string   `json:"parameters"` // Parameters will be JSON serialized to provide flexibility
}
*/

$("#submitcreationvps").on('click',function(){


  //Custom HTML-JSON Serialization for VPS
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

  

  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"createVPS",time:"",status:"",result:"",parameters:"["+JSON.stringify(createVpsJSON)+"]"};

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


/*This Job will respect the following JSON Structure on "parameters":
type Domain struct {
    Name string   `json:"name"`
    Active   string `json:"active"`         // It is being used by an Implant or not    
    Dtype string `json:"dtype"`             // Godaddy,Facebook,...
    Domain string   `json:"domain"`         // Just for domain providers
    Parameters string   `json:"parameters"`
}
*/
$("#params").on('click','#submitcreationdomain',function(){


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

//Different Creation form for SaaS
$("#params").on('click','#submitcreationdomainSaaS',function(){


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
      break; 
    case 'ssh_rev_shell':
      $("#params").empty();
      $("#params").append(`
    <form role="form" id="stagingparamsform">
    </form> 
      `); 
      break; 
  }
});

/*This Job will respect the following JSON Structure on "parameters":
type Staging struct {
    Name string   `json:"name"`
    Stype string  `json:"stype"`
    TunnelPort string  `json:"tunnelport"`            // Interactive stager, dropplet, Tunneler...
    Parameters string   `json:"parameters"` // Parameters will be JSON serialized to provide flexibility
    VpsName        string   `json:"vpsname"`
    DomainName        string   `json:"domainame"`
}
*/

$("#submitcreationstaging").on('click',function(){


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
          if (response != null){
            return
          }
        }

    });

});

/*This Job will respect the following JSON Structure on "parameters":
type Report struct {
    Name string   `json:"name"`
}
*/


$("#submitcreationreport").on('click',function(){

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

/* This Job will respect the following JSON Structure on "parameters":
type AddOperator struct {
    Username string   `json:"username"`
    Password string   `json:"password"`
}
*/


$("#submitaddOperator").on('click',function(){

  function objectifyForm(formArrayVps) {
    var returnArray = {};
    for (var i = 0; i < formArrayVps.length; i++){
        returnArray[formArrayVps[i]['name']] = formArrayVps[i]['value'];
    }
    return returnArray;
  }

  //Serialize form in the correct way
  var addOperatorJSON = objectifyForm($("#addoperatorform").serializeArray());

  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"addOperator",time:"",status:"",result:"",parameters:"["+JSON.stringify(addOperatorJSON)+"]"};

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




