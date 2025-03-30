export async function atLeastOnePromise<T>(promises: Promise<T>[]): Promise<[T[], [number, Error][]]> {
    const results = await Promise.allSettled(promises);

    let fulfilled: T[] = [];
    let rejected: [number, Error][] = [];

    for (const [index, result] of results.entries()) {
        if (result.status === 'fulfilled') fulfilled.push(result.value);
        else rejected.push([index, result.reason]);
    };

    if (fulfilled.length === 0 && results.length > 0) {
        throw "All promises failed";
    }

    return [fulfilled, rejected];
}

export async function deepCopy<T>(obj: T): Promise<T> {
    if (obj === null || obj === undefined) return obj;

    if (isPrimitive(obj)) return obj;

    if (obj instanceof Date) return new Date(obj.getTime()) as T;

    if (obj instanceof Array) {
        const arr = obj as Array<any>;
        const copy = new Array(arr.length);
        for (let i = 0; i < arr.length; i++) {
            copy[i] = await deepCopy(arr[i]);
        }
        return copy as T;
    }

    if (obj instanceof Object) {
        const copy = {} as {[key: string]: any };
        for (let key in obj) {
            copy[key] = await deepCopy((obj as any)[key]);
        }
        return copy as T;
    }

    return JSON.parse(JSON.stringify(obj)) as T;
}

// https://stackoverflow.com/questions/25456013/javascript-deepequal-comparison
export function deepEquality<T>(a: T, b: T): boolean {
    if (a === b)
        return true;

    if (isPrimitive(a) && isPrimitive(b))
        return a === b;

    if (a instanceof Date || b instanceof Date) {
        return JSON.stringify(a) === JSON.stringify(b);
    }

    if (Object.keys(a as Object).length !== Object.keys(b as Object).length)
        return false;

    for (let key in a) {
        if(!(key in (b as Object))) return false;
        if(!deepEquality(a[key], b[key])) return false;
    }

    return true;
}

function isPrimitive(obj: any) {
    return (obj !== Object(obj));
}

export function isDescendentOf(descendent: HTMLElement, element: HTMLElement): boolean {
    for (let node: (HTMLElement | null) = descendent; node; node = node.parentElement) {
        if (node === element) return true;
    }
    return false;
}

export function isChildOfModal(element: HTMLElement): boolean {
    for (let node: (HTMLElement | null) = element; node; node = node.parentElement) {
        if (node instanceof HTMLDialogElement) return true;
    }
    return false;
}

export function isChildOfForm(element: HTMLElement): boolean {
    for (let node: (HTMLElement | null) = element; node; node = node.parentElement) {
        if (node instanceof HTMLFormElement) return true;
    }
    return false;
}