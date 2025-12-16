Output assets/demo.gif

Set Theme "Catppuccin Mocha"
Set FontSize 14
Set Padding 20
Set Margin 20
Set MarginFill "#1e1e2e"
Set BorderRadius 10
Set Framerate 30


# Step 1: Backup in same directory
Type echo "Option 1: Backup in same directory"
Enter
Sleep 1s
Enter
Type "gb example.txt"
Enter
Sleep 1s
Type "ls -lh example.txt*"
Enter
Sleep 2s
Enter
Type "clear"
Enter
Sleep 500ms

# Step 2: Backup to target directory
Type echo "Option 2: Backup to specific target directory"
Enter
Sleep 1s
Enter
Type "mkdir -p backup-target"
Enter
Sleep 300ms
Type "gb example.txt to backup-target"
Enter
Sleep 1s
Type "ls -lh backup-target/"
Enter
Sleep 2s
Enter
Type "clear"
Enter
Sleep 500ms

# Step 3: Backup in current directory
Type echo "Step 3: Backup in current directory"
Enter
Sleep 1s
Enter
Type "cd backup-target"
Enter
Sleep 300ms
Type "gb ../example.txt ."
Enter
Sleep 1s
Type "ls -lh"
Enter
Sleep 2s
Enter
Type "clear"
Enter
Sleep 500ms

# Cleanup
Hide
Type "rm -rf example.txt* backup-target"
Enter
Sleep 500ms
