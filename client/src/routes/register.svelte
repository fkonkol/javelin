<svelte:head>
	<title>Sign Up</title>
</svelte:head>

<main class="mx-auto flex justify-center items-center w-screen h-screen font-sans">
	<article class="shadow-xl w-full max-w-md p-8 rounded-xl">
		<svg
			class="mx-auto h-14 mb-3 pointer-events-none select-none"
			xmlns="http://www.w3.org/2000/svg"
			viewBox="0 0 224.42 151.34"
			><defs
				><style>
					.cls-1,
					.cls-2 {
						fill: #333;
					}
					.cls-2 {
						font-size: 60px;
						font-family: Gilroy-ExtraBold, Gilroy;
						font-weight: 800;
					}
					.cls-3 {
						letter-spacing: -0.04em;
					}
					.cls-4 {
						letter-spacing: -0.02em;
					}
					.cls-5 {
						fill: none;
					}
				</style></defs
			><g id="Layer_2" data-name="Layer 2"
				><g id="Layer_1-2" data-name="Layer 1"
					><path
						class="cls-1"
						d="M40.33,93.82h0V69.89a8.26,8.26,0,0,0-8.25-8.25H8.25A8.26,8.26,0,0,0,0,69.89V93.71A8.25,8.25,0,0,0,8.25,102H32l4.62,2.76a2.42,2.42,0,0,0,1.29.36,2.54,2.54,0,0,0,2.54-2.57Zm-23-12.2a4.13,4.13,0,1,1,4.12-4.13A4.13,4.13,0,0,1,17.37,81.62Zm13,0a4.13,4.13,0,1,1,4.12-4.13A4.13,4.13,0,0,1,30.37,81.62Z"
					/><text /><text /><text class="cls-2" transform="translate(50.37 102.07)"
						><tspan class="cls-3">j</tspan><tspan class="cls-4" x="14.1" y="0">a</tspan><tspan
							class="cls-3"
							x="50.28"
							y="0">veli</tspan
						><tspan x="139.08" y="0">n</tspan></text
					><path
						class="cls-5"
						d="M34.49,77.49a4.13,4.13,0,1,1-4.12-4.12A4.12,4.12,0,0,1,34.49,77.49Z"
					/><path
						class="cls-5"
						d="M21.49,77.49a4.13,4.13,0,1,1-4.12-4.12A4.12,4.12,0,0,1,21.49,77.49Z"
					/></g
				></g
			></svg
		>

		<p class="text-center font-extralight text-2xl mb-2">Welcome to Javelin!</p>

		<p class="text-center text-sm mb-8">
			Already have an account?
			<a class="underline" href="/login">Sign in</a>
		</p>

		<form>
      <div class="w-full h-11 rounded border-gray-100 border-2 my-5">
				<input
					class="w-full h-full bg-transparent px-4 outline-none"
					name="email"
					type="email"
					placeholder="Email"
				/>
			</div>

			<div class="w-full h-11 rounded border-gray-100 border-2 my-5">
				<input
					class="w-full h-full bg-transparent px-4 outline-none"
					name="username"
					type="text"
					placeholder="Username"
				/>
			</div>

			<div class="w-full h-11 rounded border-gray-100 border-2 my-5">
				<input
					class="w-full h-full bg-transparent px-4 outline-none"
					name="password"
					type="password"
					placeholder="Password"
				/>
			</div>

			<input
				class="bg-black text-white w-full h-11 border-none border-r-2 uppercase font-medium tracking-wide text-sm hover:cursor-pointer rounded"
				type="submit"
				value="Sign Up"
			/>
		</form>
	</article>
</main>
