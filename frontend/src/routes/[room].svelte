<script context="module">
	/**
	 * @type {import('@sveltejs/kit').Load}
	 */
	export async function load({ page, fetch }) {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/publicRoom/${page.params['room']}`);

		if (res.ok) {
			return {
				props: {
					room: await res.json()
				}
			};
		}

		return {
			status: res.status,
			error: new Error(`Could not load that room`)
		};
	}
</script>

<script lang="ts">
	import type { PublicRoom, Option } from '$lib/types/Room';

	export let room: PublicRoom;

	let selectedOption: string | undefined = undefined;

	let error = '';
	let loading = false;

	$: hasSelectedOptionAlready = room.options.find((opt) => opt.selectedByMe === true);

	const submitOption = async (optionID: string | undefined, roomID: string) => {
		if (optionID === undefined) {
			error = 'Please select an option';
			return;
		}

		loading = true;
		error = '';

		await fetch(`${import.meta.env.VITE_API_URL}/room/${roomID}/option/${optionID}/select`, {
			method: 'PATCH',
			headers: {
				accepts: 'application/json'
			}
		})
			.then((res) => res.json())
			.then((res: Option) => {
				const updatedOption = room.options.find(({ id }) => id === res.id);
				const nonUpdatedOptions = room.options.filter(({ id }) => id !== res.id);

				if (updatedOption) {
					room = {
						...room,
						options: [...nonUpdatedOptions, { ...updatedOption, selectedByMe: true }].sort((a, b) =>
							a.value.localeCompare(b.value)
						)
					};
				}
			})
			.catch(() => {
				error = 'That option may have already been selected, try refreshing the page.';
			})
			.finally(() => {
				loading = false;
			});
	};
</script>

<div class="container mx-auto max-w-lg py-4">
	<form
		on:submit|preventDefault={() => submitOption(selectedOption, room.id)}
		class="card shadow-lg"
	>
		<div class="card-body">
			<div class="card-title">
				{room.question}
			</div>
			<div class="flex flex-col space-y-2">
				{#each room.options.sort((a, b) => a.value.localeCompare(b.value)) as option}
					<button
						class:btn-disabled={!option.available || !!hasSelectedOptionAlready}
						disabled={!option.available || !!hasSelectedOptionAlready}
						aria-disabled={!option.available || !!hasSelectedOptionAlready}
						class:btn-active={selectedOption === option.id ||
							hasSelectedOptionAlready?.id === option.id}
						class:white={selectedOption === option.id || hasSelectedOptionAlready?.id === option.id}
						on:click={() => (selectedOption = option.id)}
						type="button"
						class="btn btn-secondary btn-outline border-2 not-uppercase"
					>
						{option.value}
					</button>
				{/each}
			</div>
			<div class="flex justify-end mt-4">
				<button
					class="btn btn-primary"
					class:btn-loading={loading}
					class:btn-disabled={hasSelectedOptionAlready}
					disabled={!!hasSelectedOptionAlready}
					aria-disabled={!!hasSelectedOptionAlready}>Save</button
				>
			</div>
			{#if error}
				<div class="alert alert-warning mt-4">{error}</div>
			{/if}
			{#if hasSelectedOptionAlready}
				<div class="alert alert-info mt-4">You have selected {hasSelectedOptionAlready.value}</div>
			{/if}
		</div>
	</form>
</div>

<style lang="postcss">
	.white {
		@apply text-white;
	}

	.not-uppercase {
		text-transform: none;
	}
</style>
