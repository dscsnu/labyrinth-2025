<script lang="ts">
    import { onMount } from 'svelte';
    import QRCode from 'qrcode';
    
    let qrDataUrl = $state('');
    
    async function generateQRCode() {
      try {
        const payload = 'location:1234';

        const encryptResponse = await fetch('/api/qr/encrypt', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ payload })
        });
        
        if (!encryptResponse.ok) {
          throw new Error('Failed to encrypt payload');
        }
        
        const encryptData = await encryptResponse.json();
    
        qrDataUrl = await QRCode.toDataURL(encryptData.token, {
          width: 300,
          margin: 2,
          color: {
            dark: '#000000',
            light: '#ffffff'
          }
        });

      } catch (err) {
        console.error('Error generating QR code:', err);
      }
    }

    onMount(() => {
      generateQRCode();
      const intervalId = window.setInterval(generateQRCode, 5000);
      return () => {
        if (intervalId) clearInterval(intervalId);
      };
    });
</script>
  
<main class="flex flex-col items-center justify-center min-h-screen p-4">
    <h1 class="text-2xl font-bold mb-4">QR Code</h1>

    <div class="border-2 p-4 rounded-lg bg-white">
            <img src={qrDataUrl} alt="QR Code" width="300" height="300" />
    </div>
</main>