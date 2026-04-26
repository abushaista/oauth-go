const { createApp, ref, onMounted } = Vue;

const SCOPE_MAP = {
    'openid': {
        icon: '🔑',
        label: 'OpenID Connect',
        description: 'Verify your identity'
    },
    'profile': {
        icon: '👤',
        label: 'Profile',
        description: 'View your basic profile info (name, avatar)'
    },
    'email': {
        icon: '✉️',
        label: 'Email',
        description: 'View your email address'
    },
    'offline_access': {
        icon: '🔄',
        label: 'Offline Access',
        description: 'Maintain access while you are not using the app'
    }
};

const ConsentApp = {
    setup() {
        const loading = ref(true);
        const isSubmitting = ref(false);
        const errorMessage = ref('');
        const clientName = ref('Application');
        const clientId = ref('');
        const redirectUri = ref('');
        const state = ref('');
        const codeChallenge = ref('');
        const codeChallengeMethod = ref('');
        const scopeStr = ref('');
        const scopes = ref([]);

        const parseParams = () => {
            const params = new URLSearchParams(window.location.search);
            clientId.value = params.get('client_id') || '';
            clientName.value = params.get('client_id') || 'Application';
            redirectUri.value = params.get('redirect_uri') || '';
            state.value = params.get('state') || '';
            scopeStr.value = params.get('scope') || 'openid profile email';
            codeChallenge.value = params.get('code_challenge') || '';
            codeChallengeMethod.value = params.get('code_challenge_method') || 'S256';

            // Parse scopes into display items
            const scopeNames = scopeStr.value.split(' ').filter(Boolean);
            scopes.value = scopeNames.map(name => {
                const mapped = SCOPE_MAP[name];
                if (mapped) {
                    return { id: name, ...mapped };
                }
                return {
                    id: name,
                    icon: '📋',
                    label: name,
                    description: `Access to ${name}`
                };
            });
        };

        const handleAllow = async () => {
            isSubmitting.value = true;
            errorMessage.value = '';

            try {
                // POST consent approval
                const response = await fetch('/consent', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        client_id: clientId.value,
                        scopes: scopeStr.value.split(' '),
                        approved: true
                    })
                });

                if (!response.ok) {
                    const result = await response.json().catch(() => ({}));
                    throw new Error(result.message || 'Failed to submit consent');
                }

                // Redirect back to authorize endpoint to complete the flow
                const authUrl = new URL('/oauth/authorize', window.location.origin);
                authUrl.searchParams.set('client_id', clientId.value);
                authUrl.searchParams.set('redirect_uri', redirectUri.value);
                authUrl.searchParams.set('response_type', 'code');
                authUrl.searchParams.set('scope', scopeStr.value);
                authUrl.searchParams.set('code_challenge', codeChallenge.value);
                authUrl.searchParams.set('code_challenge_method', codeChallengeMethod.value);
                if (state.value) {
                    authUrl.searchParams.set('state', state.value);
                }

                window.location.href = authUrl.toString();
            } catch (err) {
                errorMessage.value = err.message;
                isSubmitting.value = false;
            }
        };

        const handleDeny = () => {
            if (!redirectUri.value) {
                window.location.href = '/ui/login.html';
                return;
            }

            const denyUrl = new URL(redirectUri.value);
            denyUrl.searchParams.set('error', 'access_denied');
            denyUrl.searchParams.set('error_description', 'The user denied consent');
            if (state.value) {
                denyUrl.searchParams.set('state', state.value);
            }
            window.location.href = denyUrl.toString();
        };

        onMounted(() => {
            parseParams();
            loading.value = false;
        });

        return {
            loading,
            isSubmitting,
            errorMessage,
            clientName,
            redirectUri,
            scopes,
            handleAllow,
            handleDeny
        };
    }
};

createApp(ConsentApp).mount('#consent-app');
