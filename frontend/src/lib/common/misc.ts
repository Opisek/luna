export async function atLeastOnePromise<T>(promises: Promise<T>[]): Promise<T[]> {
    const results = await Promise.allSettled(promises);

    const fulfilled = results.filter(result => result.status === 'fulfilled') as PromiseFulfilledResult<T>[];
    if (fulfilled.length === 0 && results.length > 0) {
        throw "All promises failed";
    }

    return fulfilled.map(result => result.value);
}