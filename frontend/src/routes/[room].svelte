<script context="module">
	/**
	 * @type {import('@sveltejs/kit').Load}
	 */
	export async function load({ page, fetch }) {
		try {
			const res = await fetch(
				`${import.meta.env.VITE_API_URL}/publicRoom/${page.params['room']}`
			).then((res) => res.json());

			return {
				props: {
					room: res
				}
			};
		} catch (e) {
			return {
				status: 404,
				error: new Error(`That room could not be found.`)
			};
		}
	}
</script>

<script lang="ts">
	import { sortOptions } from '$lib/helpers/sortOptions';

	import Cross from '$lib/icons/cross.svelte';

	import type { PublicRoom, PublicOption } from '$lib/types/Room';

	export let room: PublicRoom;

	let selectedOption: string | undefined = undefined;

	let name = '';

	let error = '';
	let loading = false;
	let unselectLoading = false;

	$: hasSelectedOptionAlready = room.options.find((opt) => !!opt.selectedByMeAs);

	$: if (hasSelectedOptionAlready) {
		name = hasSelectedOptionAlready.selectedByMeAs || name;
	}

	const submitOption = async (optionID: string | undefined, roomID: string, name: string) => {
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
			},
			body: JSON.stringify({ name })
		})
			.then((res) => res.json())
			.then((res: PublicOption) => {
				const updatedOption = room.options.find(({ id }) => id === res.id);
				const nonUpdatedOptions = room.options.filter(({ id }) => id !== res.id);

				if (updatedOption) {
					room = {
						...room,
						options: sortOptions([...nonUpdatedOptions, res])
					};
				}
			})
			.catch(() => {
				error = 'That option may have already been selected, try refreshing the page.';
			})
			.finally(() => {
				selectedOption = undefined;
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
			.then((res: PublicOption) => {
				const updatedOption = room.options.find(({ id }) => id === res.id);
				const nonUpdatedOptions = room.options.filter(({ id }) => id !== res.id);

				if (updatedOption) {
					room = {
						...room,
						options: sortOptions([...nonUpdatedOptions, res])
					};
				}
			})
			.catch(() => {
				error = 'Something went wrong, try refreshing the page.';
			})
			.finally(() => {
				selectedOption = undefined;
				unselectLoading = false;
			});
	};
</script>

<div class="container mx-auto max-w-lg py-4">
	<form
		on:submit|preventDefault={() => submitOption(selectedOption, room.id, name)}
		class="card bg-white shadow-lg"
	>
		<div class="card-body">
			<div class="card-title">
				{room.question}
			</div>
			<div class="flex flex-col space-y-2">
				{#each sortOptions(room.options) as option}
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
			<div class="form-control my-2">
				<label for="name" class="label">
					<span class="label-text">Name</span>
				</label>
				<input
					type="text"
					id="name"
					bind:value={name}
					required
					aria-required
					class="input input-bordered w-full"
					placeholder="Something to identify you"
				/>
			</div>
			<div class="flex justify-end mt-4">
				<button
					class="btn btn-primary"
					class:loading
					class:btn-disabled={hasSelectedOptionAlready}
					disabled={!!hasSelectedOptionAlready}
					aria-disabled={!!hasSelectedOptionAlready}>Submit</button
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
