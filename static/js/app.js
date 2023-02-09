
    function Promt() {
        let toast = function (c) {
          var {
            msg = "qweqwe",
            icon = "success",
            position = "top-end",
          } = c
          const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
              toast.addEventListener('mouseenter', Swal.stopTimer)
              toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
          })
  
          Toast.fire({})
        }
        let success = function (c) {
          const {
            msg = "",
            title = "",
            footer = "",
          } = c
  
          Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer,
          })
        }
        let error = function (c) {
          const {
            msg = "",
            title = "",
            footer = "",
          } = c
  
          Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer,
          })
        }
        async function customer(c) {
          const {
            msg = "",
            title = "",
            showConfirmButton=true,
          } = c
          const { value: formValues } = await Swal.fire({
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancleButton: true,
            showConfirmButton:showConfirmButton,
            willOpen: () => {
              if (c.willOpen !== undefined) {
                c.willOpen()
              }
            },
            preConfirm: () => {
              return [
                document.getElementById('start').value,
                document.getElementById('end').value
              ]
            },
            didOpen: () => {
              if (c.didOpen !== undefined) {
                c.didOpen()
              }
            }
          })
  
          if (formValues) {
            if (formValues.dismiss !== Swal.DismissReason.cancel) {
              if (formValues.value !== "") {
                if (c.callback !== undefined) {
                  c.callback(formValues)
                }
              } else {
                c.callback(false)
              }
            } else {
              c.callback(false)
            }
          }
  
        }
        return {
          toast: toast,
          success: success,
          error: error,
          customer: customer,
        }
      }