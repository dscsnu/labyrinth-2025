<script lang="ts">
    import { onMount } from 'svelte';
    import QRCode from 'qrcode';
    
    let qrDataUrl = $state('');
    let selectedType = $state('location');
    let codeInput = $state('1234');
    
    let payload = $state(`${selectedType}:${codeInput}`);

    async function generateQRCode() {
      try {
        payload = `${selectedType}:${codeInput}`;

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
    
      if (payload) {
        generateQRCode();
      }

    onMount(() => {
      const intervalId = window.setInterval(generateQRCode, 5000);
      return () => {
        if (intervalId) clearInterval(intervalId);
      };
    });
</script>
  
<main class="flex flex-col items-center justify-center min-h-screen p-4">
    <h1 class="text-2xl font-bold mb-4">QR Code</h1>
    
    <div class="mb-4 w-full max-w-md">
      <p class="text-sm font-medium text-gray-700 mb-2">Select Type:</p>
      <div class="flex gap-4">
        <label class="inline-flex items-center">
          <input 
            type="radio" 
            bind:group={selectedType} 
            value="location"
            class="form-radio h-4 w-4 text-blue-600"
          />
          <span class="ml-2">Location</span>
        </label>
        <label class="inline-flex items-center">
          <input 
            type="radio" 
            bind:group={selectedType} 
            value="qte"
            class="form-radio h-4 w-4 text-blue-600"
          />
          <span class="ml-2">QTE</span>
        </label>
      </div>
    </div>
    
    <div class="mb-6 w-full max-w-md">
      <label for="code-input" class="block text-sm font-medium text-gray-700 mb-2">
        Enter {selectedType === 'location' ? 'Location' : 'QTE'} Code:
      </label>
      <input
        id="code-input"
        type="text"
        bind:value={codeInput}
        placeholder="Enter code"
        class="w-full p-2 border border-gray-300 rounded-md"
      />
    </div>

    <div class="border-2 p-4 rounded-lg bg-white mb-4">
      {#if qrDataUrl}
        <img src={qrDataUrl} alt="QR Code" width="300" height="300" />
      {:else}
        <div class="w-[300px] h-[300px] flex items-center justify-center">
          <p>Loading...</p>
        </div>
      {/if}
    </div>
</main>