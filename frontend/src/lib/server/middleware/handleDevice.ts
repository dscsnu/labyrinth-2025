import { TOKEN_NAME as DEVICE_TOKEN } from "$lib/stores/DeviceStore";
import type { Handle } from "@sveltejs/kit";

const handleDevice: Handle = async({ event, resolve }) => {
    const response = await resolve(event, {
        transformPageChunk: ({ html }) => {
            const maxAge = 365 * 24 * 60 * 60;

            let currentDevice = event.cookies.get(DEVICE_TOKEN);
            if (!currentDevice) {
                const userOnMobile = event.request.headers.get('sec-ch-ua-mobile') === '?1';
                currentDevice = userOnMobile ? 'mobile' : 'desktop';

                event.cookies.set(DEVICE_TOKEN, currentDevice, {
                    path: '/',
                    expires: new Date(Date.now() + maxAge),
                    maxAge,
                    httpOnly: false,
                    sameSite: 'strict',
                })
            }

            return html
                .replace('data-device=""', `data-device="${currentDevice}"`)
        }
    })

    response.headers.set('Accept-CH', 'Sec-CH-UA-Mobile');
    response.headers.set('Vary', 'Sec-CH-UA-Mobile');
    response.headers.set('Critical-CH', 'Sec-CH-UA-Mobile');

    return response;
}

export default handleDevice;