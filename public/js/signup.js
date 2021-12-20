(function($) {
    'use strict';
    $(function() {

        
        $('#verifyEmailBtn').on("click", function(event) {
            var inputEmail = $('#email').val();
            var emailDiv = $('#emailDiv')
            var verifyEmailBtn = $('#verifyEmailBtn')
            if(inputEmail==""){
                event.preventDefault();
                return;
            }
            var data = {email: inputEmail};
            fetch("/mail/verifying", {
                method: 'POST', // or 'PUT'
                body: JSON.stringify(data), // data can be `string` or {object}!
                
              }).then(res=>res.json())
              .then(function(response){
                if(response.success){
                    verifyEmailBtn.val('이메일 인증 중')
                    verifyEmailBtn.attr('disabled', true);
                    $('#email').attr('disabled', true);
                    emailDiv.after('<div class="row button" id="checkCodeDiv"><input id="checkCodeBtn" type="button" value="확인"></div>')
                    emailDiv.after('<div class="row" id="verifyDiv"><i class="fas fa-key"></i><input id="verifyCode" type="text"></div>')
                    $('#signUpBtn').attr('disabled', false);
                }
              })

        })
        //localhost:8080/auth/code/577bfbc?email=1@naver.com
        $(document).on("click", "#checkCodeBtn", function() {
            var verifyCode = $('#verifyCode').val();
            var inputEmail = $('#email').val();
            var verifyEmailBtn = $('#verifyEmailBtn')
            fetch("/auth/code/"+verifyCode+"?email="+inputEmail, {
                method: 'GET', // or 'PUT'
              }).then(res=>res.json())
              .then(function(response){
                  if(response.success){
                    $('#checkCodeDiv').remove()
                    $('#verifyDiv').remove()
                    verifyEmailBtn.val('이메일 인증 완료')
                  }else{
                      alert('코드값이 틀렸거나 기간이 지났습니다.')
                  }
            });
       
        })
        $('#signUpBtn').on("click", function(event) {
            var inputEmail = $('#email').val();
            var inputPwd = $('#pwd').val();
            var inputCheckPwd = $('#pwd_check').val();
            var inputName = $('#name').val();
            var inputPhone = $('#phone').val();
            if(inputEmail=="" ||inputPwd=="" ||inputCheckPwd=="" ||inputName=="" ||inputPhone=="" ){
                alert("모든 값을 입력하시오")
                return
            }
            if(inputCheckPwd!=inputPwd){
                alert("값 다름")
                return
            }
            $.post("/api/signup", {email:inputEmail, password:inputPwd, name:inputName, phone:inputPhone}, test);
            })
        
            var test = function(res){
                if(res.error){
                    alert(res.error)
                    return
                }
                window.location.href = '/'
            }
    
    });
    })(jQuery);