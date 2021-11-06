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
	import Cross from '$lib/icons/cross.svelte';

	import type { PublicRoom, Option } from '$lib/types/Room';

	export let room: PublicRoom;

	let selectedOption: string | undefined = undefined;

	let error = '';
	let loading = false;
	let unselectLoading = false;

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
						options: [
							...nonUpdatedOptions,
							{ ...updatedOption, selectedByMe: true, available: false }
						].sort((a, b) => a.value.localeCompare(b.value))
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

	const unselectOption = async (optionID: string, roomID: string) => {
		unselectLoading = true;
		error = '';

		await fetch(`${import.meta.env.VITE_API_URL}/room/${roomID}/option/${optionID}/unselect`, {
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
						options: [
							...nonUpdatedOptions,
							{ ...updatedOption, selectedByMe: false, available: true }
						].sort((a, b) => a.value.localeCompare(b.value))
					};
				}
			})
			.catch(() => {
				error = 'Something went wrong, try refreshing the page.';
			})
			.finally(() => {
				unselectLoading = false;
			});
	};
</script>

<div class="container mx-auto max-w-lg py-4">
	<form
		on:submit|preventDefault={() => submitOption(selectedOption, room.id)}
		class="card bg-white shadow-lg"
	>
		<div class="card-body">
			<div class="card-title">
				{room.question}
			</div>
			<div class="flex flex-col space-y-2">
				{#each room.options.sort((a, b) => a.value.localeCompare(b.value)) as option}
					{#if !option.available || !!hasSelectedOptionAlready}
						<div class="flex space-x-2">
							{#if hasSelectedOptionAlready?.id === option.id}
								<div class="bg-secondary text-secondary-content p-3 rounded-xl w-full">
									{option.value}
								</div>
								<button
									class:loading={unselectLoading}
									on:click={() => unselectOption(option.id, room.id)}
									type="button"
									class="btn btn-accent btn-circle"
								>
									{#if !unselectLoading}
										<Cross />
									{/if}
								</button>
							{:else}
								<div class="bg-gray-200 p-3 rounded-xl w-full">
									{option.value}
								</div>
							{/if}
						</div>
					{:else}
						<button
							on:click={() => (selectedOption = option.id)}
							type="button"
							class:btn-active={selectedOption === option.id}
							class:white={selectedOption === option.id}
							class="btn btn-block btn-secondary btn-outline border-2 not-uppercase"
						>
							{option.value}
						</button>
					{/if}
				{/each}
			</div>
			<div class="flex justify-end mt-4">
				<button
					class="btn btn-primary"
					class:loading
					class:btn-disabled={hasSelectedOptionAlready}
					disabled={!!hasSelectedOptionAlready}
					aria-disabled={!!hasSelectedOptionAlready}>Save</button
				>
			</div>
			{#if error}
				<div class="alert alert-warning mt-4">{error}</div>
			{/if}
			{#if !hasSelectedOptionAlready && room.options.filter((opt) => opt.available).length === 0}
				<div class="alert alert-warning mt-4">All the options have already been selected.</div>
			{/if}
			{#if hasSelectedOptionAlready}
				<div class="alert alert-info mt-4 justify-center">
					You have selected {hasSelectedOptionAlready.value}
				</div>
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
