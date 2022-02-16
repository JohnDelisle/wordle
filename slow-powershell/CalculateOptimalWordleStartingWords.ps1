$ErrorActionPreference = "inquire" 


# SCREW THIS - Powershell is way too slow, or I'm just too crappy a coder lol.. Moved to Go instead

function Initialize-Letters {
    # build our letters
    $letters = @{}
    foreach ($letter in @('a','b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z')) {
        $letters[$letter] = 0
    } 
    
    return $letters
}

function Initialize-Words {
    # import the filtered word list (see the shell script)
    $words = @{}
    foreach ($word in Get-Content -Path C:\temp\wordle.txt <# | Get-Random -Count 10 #>) {
        
        if (($word.ToCharArray() | Sort-Object -Unique).count -ne 5) {
            # word re-uses letters, that's wasteful for a starting word
            continue 
        }

        $words[$word] = 0
    }
    return $words
}

function Measure-LetterScores($letters, $words) {
    # score each of the letters
    # score is the count of words that contain the letter

    foreach ($letter in $letters.Keys.Clone()) { 
        $letters[$letter] = ($words.Keys -like "*$letter*").count
    }

    return $letters
}

function Measure-WordScores($scoredLetters, $words) {
    # score each of the words

    foreach ($word in $words.Keys.Clone()) {

        # word score is the sum of the score of each unique letter in the word
        $words[$word] = 0
        $word.ToCharArray() | Sort-Object -Unique | ForEach-Object {  
            $words[$word] += $scoredLetters["$_"]
        }
    }
    
    return $words
}

function Test-IsDuplicatePair($wordPair, $wordPairs) {
    # the pair we're lookin for
    $testWordPair = $wordPair

    foreach ($wordPair in $wordPairs) {
        if (($wordPair.Words.Name -contains $testWordPair.Words[0].Name) -and ($wordPair.Words.Name -contains $testWordPair.Words[1].Name)) {
            return $true
        }
    }
    return $false
}



function New-WordPairs ($scoredWords) {
    # word1,word2 = score
    $wordPairs = @{}

    $shortScoredWords = $scoredWords.Keys.Clone()

    foreach($wordOne in $scoredWords.Keys) {
        # shrinking list of words to create pairs from (eliminate "cat,dog" and "dog,cat" duplication)
        $shortScoredWords = $shortScoredWords | Where-Object { $_ -ne $wordOne }
        Write-Host -NoNewline "$($shortScoredWords.count) -- "

        # for the current value of wordOne, we want to remove word-pairs where letters would be duplicated (that's wasteful for wordle first word guessing)
        $extraShortScoredWords = $shortScoredWords | Where-Object {
            ($_.ToCharArray() -notcontains ($wordOne.ToCharArray())[0]) -and 
            ($_.ToCharArray() -notcontains ($wordOne.ToCharArray())[1]) -and
            ($_.ToCharArray() -notcontains ($wordOne.ToCharArray())[2]) -and
            ($_.ToCharArray() -notcontains ($wordOne.ToCharArray())[3]) -and
            ($_.ToCharArray() -notcontains ($wordOne.ToCharArray())[4])
        }

        Write-Host "$($extraShortScoredWords.count)"
        
        foreach($wordTwo in $extraShortScoredWords) {
            $wordPairs["$wordOne,$wordTwo"] = $scoredWords["$wordOne"] + $scoredWords["$wordTwo"]
        }
    }
    Write-Host "!"

    $wordPairs
}



Write-Output "Initializing Letters.."
$letters = Initialize-Letters
Write-Output "Initialized $($letters.count) Letters"

Write-Output "Initializing Words.."
$words = Initialize-Words
Write-Output "Initialized $($words.count) Words"

Write-Output "Scoring Letters.."
$scoredLetters = Measure-LetterScores -letters $letters -words $words
Write-Output "Scored $($scoredLetters.count) Letters"

Write-Output "Scoring Words.."
$scoredWords = Measure-WordScores -scoredLetters $scoredLetters -words $words 
Write-Output "Scored $($scoredWords.count) Words"


Write-Output "Generating Word Pairs.."
$wordPairs = New-WordPairs -scoredWords $scoredWords

$wordPairs | ConvertTo-Json -Depth 99 | Out-File c:\temp\wordPairs.json