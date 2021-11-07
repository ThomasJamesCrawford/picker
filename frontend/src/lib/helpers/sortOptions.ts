import type { PublicOption } from '$lib/types/Room';

export const sortOptions = (options: PublicOption[]): PublicOption[] =>
	options.sort((a, b) => {
		if (a.value === b.value) {
			return a.id.localeCompare(b.id);
		}

		return a.value.localeCompare(b.value);
	});
