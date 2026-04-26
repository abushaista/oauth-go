const { createApp, ref, onMounted } = Vue;

const App = {
    setup() {
        const username = ref('');
        const password = ref('');
        const isLoading = ref(false);
        const errorMessage = ref('');
        const successMessage = ref('');
        const shakeError = ref(false);

        const triggerShake = () => {
            shakeError.value = true;
            setTimeout(() => {
                shakeError.value = false;
            }, 400); // Wait for the animation duration
        };

        const handleLogin = async () => {
            isLoading.value = true;
            errorMessage.value = '';
            successMessage.value = '';

            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        username: username.value,
                        password: password.value
                    })
                });

                const result = await response.json();

                if (response.ok && result.success) {
                    successMessage.value = 'Login successful! Redirecting...';
                    
                    const urlParams = new URLSearchParams(window.location.search);
                    const returnTo = urlParams.get('return_to');

                    setTimeout(() => {
                        if (returnTo) {
                            window.location.href = decodeURIComponent(returnTo);
                        } else {
                            window.location.href = '/ui/login.html?success=true';
                        }
                    }, 1500);
                } else {
                    throw new Error(result.message || 'Login failed');
                }
            } catch (error) {
                errorMessage.value = error.message;
                isLoading.value = false;
                triggerShake();
            }
        };

        onMounted(() => {
            const urlParams = new URLSearchParams(window.location.search);
            if (urlParams.get('success') === 'true') {
                successMessage.value = 'Welcome! You have been securely logged in.';
                const newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
                window.history.pushState({path:newUrl},'',newUrl);
            }
        });

        return {
            username,
            password,
            isLoading,
            errorMessage,
            successMessage,
            shakeError,
            handleLogin
        };
    }
};

createApp(App).mount('#app');
