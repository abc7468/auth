(function($) {
    'use strict';
    $(function() {
        
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