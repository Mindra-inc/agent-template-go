#!/bin/bash

# Deploy Go Template to GitHub
# This script initializes a git repo and pushes to GitHub

set -e

echo "🚀 Deploying Go Template to GitHub..."

# Initialize git if not already done
if [ ! -d ".git" ]; then
    echo "📦 Initializing git repository..."
    git init
    git branch -M main
else
    echo "✅ Git repository already initialized"
fi

# Create .gitignore if it doesn't exist
if [ ! -f ".gitignore" ]; then
    echo "📝 Creating .gitignore..."
    cat > .gitignore << 'EOF'
# Binaries
agent
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output
*.out

# Go workspace
go.work

# Environment
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
EOF
else
    echo "✅ .gitignore already exists"
fi

# Add all files
echo "📦 Adding files..."
git add .

# Commit
echo "💾 Creating commit..."
git commit -m "Initial commit: Go agent template" || echo "No changes to commit"

# Add remote if not exists
if ! git remote | grep -q "origin"; then
    echo "🔗 Adding remote repository..."
    git remote add origin https://github.com/Mindra-inc/agent-template-go.git
else
    echo "✅ Remote already exists"
fi

# Push to GitHub
echo "⬆️  Pushing to GitHub..."
echo ""
echo "⚠️  Make sure you've created the repository on GitHub:"
echo "   https://github.com/Mindra-inc/agent-template-go"
echo ""
echo "Press Enter to continue or Ctrl+C to cancel..."
read

git push -u origin main

echo ""
echo "✅ Successfully deployed to GitHub!"
echo "🌐 Repository: https://github.com/Mindra-inc/agent-template-go"
echo ""
echo "Next steps:"
echo "1. Visit the repository on GitHub"
echo "2. Add a description and topics"
echo "3. Enable Issues and Discussions (optional)"
echo "4. Add repository secrets for CI/CD (if needed)"
