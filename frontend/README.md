# Another Frontend :(

### fetchWithAuth
Automatically append jwt to outgoing request.
```ts
import { fetchWithAuth } from '$lib/utils/fetchWithAuth.js';

const testFetch = async () => {
    const response = await fetchWithAuth('https://github.com/dscsnu/labyrinth-2025');
    const data = await response.json();

    /* ---snip--- */
}
```