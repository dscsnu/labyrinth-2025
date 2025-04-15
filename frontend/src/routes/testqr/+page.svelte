<script>
    // @ts-ignore
    import { Html5Qrcode } from 'html5-qrcode';
    import { onMount, onDestroy } from 'svelte';
    import { goto } from '$app/navigation';

    let html5Qrcode;
    let scanning = false;

    onMount(() => {
        init();
    });
    
    onDestroy(() => {
        if (html5Qrcode && scanning) {
            html5Qrcode.stop().catch(console.error);
        }
    });

    async function init() {
        try {
            html5Qrcode = new Html5Qrcode('reader');
            await startScanning();
        } catch (err) {
            console.error('Error initializing scanner:', err);
            alert('Failed to initialize QR scanner');
        }
    }

    async function startScanning() {
        try {
            await html5Qrcode.start(
                { facingMode: 'environment' },
                {
                    fps: 10,
                    qrbox: { width: 250, height: 250 },
                },
                onScanSuccess,
                onScanFailure
            );
            scanning = true;
        } catch (err) {
            console.error('Error starting scanner:', err);
            alert('Failed to start scanner. Please ensure camera permissions are granted.');
        }
    }

    async function onScanSuccess(decodedText, decodedResult) {
        await html5Qrcode.stop();
        scanning = false;
        
        await decryptQrCode(decodedText);
    }

    function onScanFailure(error) {
        console.warn(`Code scan error = ${error}`);
    }
    
    async function decryptQrCode(token) {
        console.log('Decrypting QR code with token:', token);
        try {
            const response = await fetch('/qr/decrypt', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ token })
            });

            if (!response.ok) {
                throw new Error(`Decryption failed with status ${JSON.stringify(response)}`);
            }
            
            const result = await response.json();

            if (result.valid && result.route) {
                const destinationUrl = `${result.route}?token=${encodeURIComponent(token)}`;
                goto(destinationUrl);
    
            } else {
                alert(`Invalid QR code or missing route: ${JSON.stringify(result)}`);
                startScanning();
            }
        } catch (err) {
            console.error('Error decrypting QR code:', err);
            alert(`Error: ${err instanceof Error ? err.message : 'Failed to decrypt QR code'}`);
            startScanning();
        }
    }
</script>

<style>
    main {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 20px;
        padding: 20px;
    }
    
    #reader {
        width: 100%;
        max-width: 500px;
        min-height: 300px;
        background-color: #f0f0f0;
        border: 1px solid #ccc;
    }
</style>

<main>
    <h1>QR Code Scanner</h1>
    <div id="reader"></div>
</main>