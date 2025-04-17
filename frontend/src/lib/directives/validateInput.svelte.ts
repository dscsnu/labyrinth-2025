import type { Action } from "svelte/action";

export enum ValidationOptions {
    ALPHABET_LOWER = 'abcdefghijklmnopqrstuvwxyz',
    ALPHABET_UPPER = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
    NUMERIC = '0123456789',
    SPACE = ' ',
    SPECIAL_CHARACTERS = ',.!?;:_<>',
    ALL = ' ,.!?;:_<>0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
}

export interface IValidationConfig {
    allowed?: string[]
    maxLength?: number,
}

export const validateInput: Action<HTMLInputElement, IValidationConfig> = (node, config) => {
    const allowedInputs = config.allowed?.join('') || ValidationOptions.ALL;
    const maxLength = config.maxLength;

    const handleInput = (event: Event) => {
        const target = event.target as HTMLInputElement;
        const currentValue = target.value;

        let validValue = currentValue
            .split('')
            .filter(char => allowedInputs.includes(char))
            .join('');

        if (maxLength !== undefined && validValue.length > maxLength)
            validValue = validValue.substring(0, maxLength);

        if (validValue !== currentValue) {
            target.value = validValue;
            node.dispatchEvent(new Event('input', { bubbles: true }));
        }
    };

    $effect(() => {
        node.addEventListener('input', handleInput);
        return () => node.removeEventListener('input', handleInput);
    })
}