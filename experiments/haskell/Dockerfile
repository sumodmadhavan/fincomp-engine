# Use the official Haskell image as the base image
FROM haskell:8.10.7

# Set the working directory in the container
WORKDIR /app

# Copy all the files from the project directory to the working directory
COPY . .

# Update cabal and install dependencies
RUN cabal update && cabal build --only-dependencies -j4

# Build the application
RUN cabal build

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["cabal", "run"]
